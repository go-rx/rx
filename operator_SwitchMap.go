package rx

func SwitchMap[T any, S any](source Observable[T], project func(T) Observable[S]) Observable[S] {
	return Func(func(subscriber Writer[S]) (err error) {
		outerWriter, outerReader := Pipe(PipeWithParentLifecycle[T](subscriber))
		source.Subscribe(outerWriter)
		var innerLifecycle Lifecycle
		for {
			var ok bool
			var outerValue T
			var innerObservable Observable[S]

			if outerValue, ok = outerReader.Read(); !ok {
				if innerLifecycle != nil {
					innerLifecycle.Wait()
				}
				return
			}

			if innerLifecycle != nil {
				innerLifecycle.Kill(nil)
			}

			innerObservable = project(outerValue)
			innerWriter, innerReader := Pipe(PipeWithParentLifecycle[S](subscriber))
			subscriber.Go(func() (err error) {
				for {
					if innerValue, ok := innerReader.Read(); ok {
						// forward inner observable value to outter subscriber
						if !subscriber.Write(innerValue) {
							return
						}
					} else {
						return
					}
				}
			})

			innerLifecycle = innerReader
			innerObservable.Subscribe(innerWriter)
		}
	})
}
