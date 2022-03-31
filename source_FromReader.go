package rx

func FromReader[T any](reader Reader[T]) Observable[T] {
	return Func(func(subscriber Writer[T]) (err error) {
		for {
			if value, ok := reader.Read(); ok {
				if !subscriber.Write(value) {
					return
				}
			} else {
				return
			}
		}
	})
}
