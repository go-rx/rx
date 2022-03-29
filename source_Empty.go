package rx

func Empty[T any]() Observable[T] {
	return Func(func(subscriber Writer[T]) (err error) {
		return
	})
}
