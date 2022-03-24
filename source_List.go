package rx

func List[T any](list []T) Observable[T] {
	return Func(func(subscriber Writer[T]) (err error) {
		for _, value := range list {
			if !subscriber.Alive() {
				return
			}
			if !subscriber.Write(value) {
				return
			}
		}
		return
	})
}
