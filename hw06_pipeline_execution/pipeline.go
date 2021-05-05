package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type (
	Stage         func(in In) (out Out)
	StageWithDone func(in In, done In) Out
)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if stages == nil {
		return nil
	}

	stagesWithDone := make([]StageWithDone, 0, len(stages))

	for _, stage := range stages {
		stagesWithDone = append(stagesWithDone, wrapStageWithDone(stage))
	}

	return combineStagesToChain(stagesWithDone)(in, done)
}

// solution with recursion, possible out of memory.
func combineStagesToChain(stages []StageWithDone) StageWithDone {
	stage := stages[0]

	if len(stages) == 1 {
		return stage
	}

	return func(in In, done In) Out {
		return combineStagesToChain(stages[1:])(stage(in, done), done)
	}
}

func wrapStageWithDone(stage Stage) func(in In, done In) Out {
	return func(in In, done In) Out {
		out := make(Bi)

		go func() {
			defer close(out)
			for value := range stage(in) {
				select {
				case <-done:
					return
				default:
					out <- value
				}
			}
		}()

		return out
	}
}
