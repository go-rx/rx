package rx

func ConcatMap[T any, S any](source Observable[T], project func(T) Observable[S]) Observable[S] {
	return MergeMap(source, project, MergeMapWithStandbyConcurrency(1), MergeMapWithMaxConcurrency(1))
}
