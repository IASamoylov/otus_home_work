package hw06pipelineexecution

import (
	"runtime/debug"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

const (
	sleepPerStage = time.Millisecond * 100
	fault         = sleepPerStage / 2
)

func TestPipeline(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					newVal := f(v)
					out <- newVal
				}
			}()
			return out
		}
	}

	stages := []Stage{
		g("Dummy", func(v interface{}) interface{} { return v }),
		g("Multiplier (* 2)", func(v interface{}) interface{} { return v.(int) * 2 }),
		g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }),
		g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }),
	}

	t.Run("simple case", func(t *testing.T) {
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		out := ExecutePipeline(in, nil, stages...)

		for s := range out {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Equal(t, []string{"102", "104", "106", "108", "110"}, result)
		require.Less(t,
			int64(elapsed),
			// ~0.8s for processing 5 values in 4 stages (100ms every) concurrently
			int64(sleepPerStage)*int64(len(stages)+len(data)-1)+int64(fault))
	})

	t.Run("done case", func(t *testing.T) {
		in := make(Bi)
		done := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		// Abort after 200ms
		abortDur := sleepPerStage * 2
		go func() {
			<-time.After(abortDur)
			close(done)
		}()

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		start := time.Now()
		for s := range ExecutePipeline(in, done, stages...) {
			result = append(result, s.(string))
		}
		elapsed := time.Since(start)

		require.Len(t, result, 0)
		require.Less(t, int64(elapsed), int64(abortDur)+int64(fault))
	})
}

func TestByMe(t *testing.T) {
	// Stage generator
	g := func(_ string, f func(v interface{}) interface{}) Stage {
		return func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					time.Sleep(sleepPerStage)
					newVal := f(v)
					out <- newVal
				}
			}()
			return out
		}
	}

	t.Run("returns in chan when stages are empty", func(t *testing.T) {
		t.Run("in is nil", func(t *testing.T) {
			require.Nil(t, ExecutePipeline(nil, nil))
		})

		t.Run("in is data", func(t *testing.T) {
			in := make(Bi)
			data := []int{1, 2, 3, 4, 5}

			go func() {
				for _, v := range data {
					in <- v
				}
				close(in)
			}()

			result := make([]int, 0, 5)
			out := ExecutePipeline(in, nil)

			for v := range out {
				result = append(result, v.(int))
			}

			require.Equal(t, data, result)
		})
	})

	t.Run("solution not throw runtime error out of memory", func(t *testing.T) {
		debug.SetMaxStack(1 << 13)

		stages := []Stage{
			g("Dummy", func(v interface{}) interface{} { return v }),
		}

		i := 0

		for i != 50 {
			stages = append(stages, g("Adder (+ 100)", func(v interface{}) interface{} { return v.(int) + 100 }))
			i++
		}

		stages = append(stages, g("Stringifier", func(v interface{}) interface{} { return strconv.Itoa(v.(int)) }))

		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		result := make([]string, 0, 10)
		for s := range ExecutePipeline(in, nil, stages...) {
			result = append(result, s.(string))
		}

		require.Equal(t, []string{"5001", "5002", "5003", "5004", "5005"}, result)
	})

	t.Run("the pipeline does not block when no one is reading from the buffered output channel", func(t *testing.T) {
		in := make(Bi, 5)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()
		result := make([]int, 0, 5)
		var wg sync.WaitGroup
		wg.Add(5)

		stages := []Stage{
			g("N * N", func(v interface{}) interface{} { return v.(int) * v.(int) }),
			g("N * N", func(v interface{}) interface{} { return v.(int) * v.(int) }),
			g("Result", func(v interface{}) interface{} {
				result = append(result, v.(int))
				wg.Done()
				return v.(int)
			}),
		}

		ExecutePipeline(in, nil, stages...)

		wg.Wait()

		require.Equal(t, []int{1, 16, 81, 256, 625}, result)
	})

	t.Run("can interrupt pipeline when no one read from out channel", func(t *testing.T) {
		defer goleak.VerifyNone(t)
		in := make(Bi)
		data := []int{1, 2, 3, 4, 5}

		go func() {
			for _, v := range data {
				in <- v
			}
			close(in)
		}()

		done := make(Bi)
		ExecutePipeline(in, done, func(in In) Out {
			return in
		})

		go func() {
			time.Sleep(time.Millisecond * 200)
			close(done)
		}()

		time.Sleep(time.Millisecond * 400)
	})
}

func TestPipelineByAlekseyBakin(t *testing.T) {
	t.Run("done", func(t *testing.T) {
		waitCh := make(chan struct{})
		defer close(waitCh)

		stageFn := func(in In) Out {
			out := make(Bi)
			go func() {
				defer close(out)
				for v := range in {
					<-waitCh
					out <- v
				}
			}()
			return out
		}

		in := make(Bi)
		const testValue = "test"
		go func() {
			in <- testValue
			close(in)
		}()

		doneCh := make(Bi)
		var resValue interface{}
		out := ExecutePipeline(in, doneCh, stageFn, stageFn, stageFn)
		close(doneCh)

		require.Eventually(t, func() bool {
			select {
			case resValue = <-out:
				return true
			default:
				return false
			}
		}, time.Second, time.Millisecond)

		require.Nil(t, resValue)
	})
}
