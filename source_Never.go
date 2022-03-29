package rx

func Never[T any]() Observable[T] {
	return Func(func(subscriber Writer[T]) (err error) {
		<-subscriber.Dying()
		return
	})
}
