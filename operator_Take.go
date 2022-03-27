package rx

func Take[T any](source Observable[T], n int) Observable[T] {
	return Func(func(subscriber Writer[T]) (err error) {
		writer, reader := Pipe[T](subscriber)
		source.Subscribe(writer)
		count := 0
		if n <= count {
			return
		}
		for {
			if input, ok := reader.Read(); ok {
				if !subscriber.Write(input) {
					return
				}
				count++
				if n <= count {
					subscriber.Kill(nil)
					return
				}
			} else {
				return
			}
		}
	})
}
