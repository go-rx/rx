package rx

type sinkT[T any] struct {
	Lifecycle
}

func Sink[T any](parent Lifecycle, options ...SinkOption) (sink Writer[T]) {
	opts := buildOptions(options)
	if opts.childLifecycle == nil {
		opts.childLifecycle = NewLifecycle()
	}
	Bind(parent, opts.childLifecycle)
	sink = &sinkT[T]{opts.childLifecycle}
	return
}

func (s *sinkT[T]) Write(value T) (ok bool) {
	return s.Alive()
}

type sinkOptions struct {
	childLifecycle Lifecycle
}

type SinkOption func(*sinkOptions)

func SinkWithChildLifecycle(child Lifecycle) SinkOption {
	return func(o *sinkOptions) {
		o.childLifecycle = child
	}
}
