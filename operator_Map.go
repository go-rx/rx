package rx

func Map[T any, S any](source Observable[T], project func(T) (S, error)) Observable[S] {
	return Func(func(subscriber Writer[S]) (err error) {
		writer, reader := Pipe(PipeWithParentLifecycle[T](subscriber))
		source.Subscribe(writer)
		for {
			if input, ok := reader.Read(); ok {
				var output S
				if output, err = project(input); err != nil {
					return
				}
				if !subscriber.Write(output) {
					return
				}
			} else {
				return
			}
		}
	})
}
