package rx

type writerT[T any] struct {
	Lifecycle
	c chan<- any
}

type readerT[T any] struct {
	Lifecycle
	c <-chan any
}

func Pipe[T any](parent Lifecycle, options ...PipeOption) (writer Writer[T], reader Reader[T]) {
	opts := buildOptions(options)
	if opts.childLifecycle == nil {
		opts.childLifecycle = NewLifecycle()
	}
	if opts.childChannel == nil {
		opts.childChannel = make(chan any)
	}
	Bind(parent, opts.childLifecycle)
	writer = &writerT[T]{opts.childLifecycle, opts.childChannel}
	reader = &readerT[T]{opts.childLifecycle, opts.childChannel}
	return
}

func (w *writerT[T]) Write(value T) (ok bool) {
	select {
	case w.c <- value:
		ok = true
	case <-w.Dying():
	}
	return
}

func (r *readerT[T]) Read() (value T, ok bool) {
	var v any
	select {
	case v, ok = <-r.c:
		value, ok = v.(T)
	case <-r.Dying():
	}
	return
}

type pipeOptions struct {
	childLifecycle Lifecycle
	childChannel   chan any
}

type PipeOption func(*pipeOptions)

func PipeWithChildLifecycle(child Lifecycle) PipeOption {
	return func(o *pipeOptions) {
		o.childLifecycle = child
	}
}

func PipeWithChildChannel(c chan any) PipeOption {
	return func(o *pipeOptions) {
		o.childChannel = c
	}
}
