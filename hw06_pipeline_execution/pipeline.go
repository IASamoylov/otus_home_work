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

	for _, stage := range stagesWithDone {
		in = stage(in, done)
	}

	return in
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
