package rx

func Of[T any](value T) Observable[T] {
	return Func(func(subscriber Writer[T]) (err error) {
		subscriber.Write(value)
		return
	})
}
