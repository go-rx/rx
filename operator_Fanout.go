package rx

func Fanout[T any](source Observable[[]T]) Observable[T] {
	return ConcatMap(source, func(list []T) Observable[T] {
		return List(list)
	})
}
