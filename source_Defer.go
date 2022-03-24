package rx

func Defer[T any](f func(Lifecycle) Observable[T]) Observable[T] {
	return Func(func(subscriber Writer[T]) (err error) {
		f(subscriber).Subscribe(subscriber)
		return
	})
}
