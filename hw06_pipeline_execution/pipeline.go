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

func wrapStageWithDone(stage Stage) StageWithDone {
	return func(in In, done In) Out {
		out := make(Bi, cap(in))

		go func() {
			defer close(out)

			stageOut := stage(in)

			for {
				select {
				case <-done:
					return
				case value, ok := <-stageOut:
					if ok {
						select {
						case <-done:
							return
						default:
							out <- value
						}
					} else {
						return
					}
				}
			}
		}()

		return out
	}
}
