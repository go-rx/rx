package rx

func Error[T any](err error) Observable[T] {
	return Func(func(subscriber Writer[T]) error {
		subscriber.Kill(err)
		return nil
	})
}
