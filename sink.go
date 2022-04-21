package rx

type sinkT[T any] struct {
	Lifecycle
}

func Sink[T any](parent Lifecycle, options ...SinkOption) (sink Writer[T]) {
	opts := buildOptions(options)
	if opts.childLifecycle == nil {
		opts.childLifecycle = NewLifecycle()
	}
	Bind(parent, opts.childLifecycle, opts.bindOptions...)
	sink = &sinkT[T]{opts.childLifecycle}
	return
}

func (s *sinkT[T]) Write(value T) (ok bool) {
	return s.Alive()
}

type sinkOptions struct {
	childLifecycle Lifecycle
	bindOptions    []BindOption
}

type SinkOption func(*sinkOptions)

func SinkWithChildLifecycle(child Lifecycle) SinkOption {
	return func(o *sinkOptions) {
		o.childLifecycle = child
	}
}

func SinkWithDiscardChildError() SinkOption {
	return func(o *sinkOptions) {
		o.bindOptions = append(o.bindOptions, BindWithDiscardChildError())
	}
}
