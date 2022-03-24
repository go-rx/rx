package rx

func Chan[T any](c <-chan T) Observable[T] {
	return Func(func(subscriber Writer[T]) (err error) {
		for {
			select {
			case value, ok := <-c:
				if !ok {
					return
				}
				if !subscriber.Write(value) {
					return
				}
			case <-subscriber.Dying():
				return
			}
		}
	})
}
