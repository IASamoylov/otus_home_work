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
		return in
	}

	for _, stage := range stages {
		in = wrapStageWithDone(stage)(in, done)
	}

	return in
}

func wrapStageWithDone(stage Stage) func(in In, done In) Out {
	return func(in In, done In) Out {
		out := make(Bi)

		go func() {
			defer close(out)
			select {
			case <-done:
				return
			default:
				for value := range stage(in) {
					select {
					case <-done:
						return
					default:
						out <- value
					}
				}
			}
		}()

		return out
	}
}
