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
		tempBuffer := make([]interface{}, 0)

		go func() {
			defer close(out)

			stageOut := stage(in)

			for {
				tempBuffer = tryToWriteBufferToOut(tempBuffer, out, done)

				select {
				case <-done:
					return
				case value, ok := <-stageOut:
					if ok {
						select {
						case <-done:
							return
						default:
							tempBuffer = writeTo(value, tempBuffer, out)
						}
					} else if len(tempBuffer) == 0 {
						return
					}
				}
			}
		}()

		return out
	}
}

func tryToWriteBufferToOut(tempBuffer []interface{}, out Bi, done In) []interface{} {
	if len(tempBuffer) == 0 {
		return tempBuffer
	}

	for {
		value := tempBuffer[0]
		select {
		case <-done:
			return tempBuffer
		case out <- value:
			tempBuffer = tempBuffer[1:]

			if len(tempBuffer) == 0 {
				return tempBuffer
			}
		default:
			return tempBuffer
		}
	}
}

func writeTo(value interface{}, tempBuffer []interface{}, out Bi) []interface{} {
	if len(tempBuffer) != 0 {
		tempBuffer = append(tempBuffer, value)
	} else {
		select {
		case out <- value: // Put 2 in the channel unless it is full
		default:
			tempBuffer = append(tempBuffer, value)
		}
	}

	return tempBuffer
}
