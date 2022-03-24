package rx

func Func[T any](f func(subscriber Writer[T]) error) Observable[T] {
	return functionObservable[T](f)
}

type functionObservable[T any] func(subscriber Writer[T]) error

func (o functionObservable[T]) Subscribe(subscriber Writer[T]) {
	subscriber.Go(func() error {
		return o(subscriber)
	})
}
