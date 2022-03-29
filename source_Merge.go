package rx

func Merge[T any](sources ...Observable[T]) Observable[T] {
	return Func(func(subscriber Writer[T]) (err error) {
		for _, source := range sources {
			source.Subscribe(subscriber)
		}
		return
	})
}
