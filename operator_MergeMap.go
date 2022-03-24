package rx

import (
	"fmt"
)

type mergeMapOptions struct {
	standbyConcurrency int
	maxConcurrency     int
}

type MergeMapOption func(*mergeMapOptions)

func MergeMap[T any, S any](source Observable[T], project func(T) Observable[S], options ...MergeMapOption) Observable[S] {
	opts := &mergeMapOptions{}
	for _, option := range options {
		option(opts)
	}
	return Func(func(subscriber Writer[S]) (err error) {
		if opts.standbyConcurrency < 0 {
			subscriber.Kill(fmt.Errorf("MergeMap: standbyConcurrency cannot be nagative: %w", ErrInvalidParameter))
			return
		}
		if opts.maxConcurrency < 0 {
			subscriber.Kill(fmt.Errorf("MergeMap: maxConcurrency cannot be nagative: %w", ErrInvalidParameter))
			return
		}
		if opts.maxConcurrency > 0 && opts.standbyConcurrency > opts.maxConcurrency {
			subscriber.Kill(fmt.Errorf("MergeMap: standbyConcurrency cannot exceed maxConcurrency: %w", ErrInvalidParameter))
			return
		}

		innerDone := make(chan bool)
		workerDone := make(chan bool)
		standbyDone := make(chan bool)
		fanoutChan := make(chan T)

		spawnWorker := func(standby bool) {
			subscriber.Go(func() (err error) {
				defer func() {
					workerDone <- true
				}()
				for {
					var ok bool
					var outerValue T
					var innerObservable Observable[S]

					if standby {
						select {
						case outerValue, ok = <-fanoutChan:
							if !ok {
								return
							}
						case <-subscriber.Dying():
							return
						case <-standbyDone:
							return
						}
					} else {
						select {
						case outerValue, ok = <-fanoutChan:
							if !ok {
								return
							}
						case <-subscriber.Dying():
							return
						}
					}

					if err = func() (err error) {
						defer func() {
							innerDone <- true
						}()

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
						innerObservable.Subscribe(innerWriter)
						return innerReader.Wait()
					}(); err != nil {
						return
					}

					if !standby {
						return
					}
				}
			})
		}

		for i := 0; i < opts.standbyConcurrency; i++ {
			spawnWorker(true)
		}

		outerWriter, outerReader := Pipe(PipeWithParentLifecycle[T](subscriber))
		subscriber.Go(func() (err error) {
			outstandingJobs := 0
			workerCount := opts.standbyConcurrency

			defer func() {
				close(fanoutChan)
				close(standbyDone)
				// drain outstandingJobs and workerCount channel
				for outstandingJobs > 0 {
					<-innerDone
					outstandingJobs--
				}
				for workerCount > 0 {
					<-workerDone
					workerCount--
				}
				// subscriber.Kill(err)
			}()

			for {
				var ok bool
				var outerValue T

				if outerValue, ok = outerReader.Read(); !ok {
					return
				}

				valueSent := false
				for !valueSent {
					if outstandingJobs == workerCount && (opts.maxConcurrency == 0 || workerCount < opts.maxConcurrency) {
						select {
						case <-workerDone:
							workerCount--
						case fanoutChan <- outerValue:
							valueSent = true
							outstandingJobs++
						case <-innerDone:
							outstandingJobs--
						case <-subscriber.Dying():
							return
						default:
							workerCount++
							spawnWorker(false)
						}
					} else {
						select {
						case <-workerDone:
							workerCount--
						case fanoutChan <- outerValue:
							valueSent = true
							outstandingJobs++
						case <-innerDone:
							outstandingJobs--
						case <-subscriber.Dying():
							return
						}
					}
				}
			}
		})

		source.Subscribe(outerWriter)
		return
	})
}

func MergeMapWithStandbyConcurrency(standbyConcurrency int) MergeMapOption {
	return func(mmo *mergeMapOptions) {
		mmo.standbyConcurrency = standbyConcurrency
	}
}

func MergeMapWithMaxConcurrency(maxConcurrency int) MergeMapOption {
	return func(mmo *mergeMapOptions) {
		mmo.maxConcurrency = maxConcurrency
	}
}
