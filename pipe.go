package rx

type writerT[T any] struct {
	Lifecycle
	c chan<- T
}

type readerT[T any] struct {
	Lifecycle
	c <-chan T
}

func Pipe[T any](options ...PipeOption[T]) (writer Writer[T], reader Reader[T]) {
	opts := &pipeOptions[T]{}
	for _, option := range options {
		option(opts)
	}
	if opts.childLifecycle == nil {
		opts.childLifecycle = NewLifecycle()
	}
	if opts.childChannel == nil {
		opts.childChannel = make(chan T)
	}
	Bind(opts.parentLifecycle, opts.childLifecycle)
	writer = &writerT[T]{opts.childLifecycle, opts.childChannel}
	reader = &readerT[T]{opts.childLifecycle, opts.childChannel}
	return
}

func (w *writerT[T]) C() chan<- T {
	return w.c
}

func (w *writerT[T]) Write(value T) (ok bool) {
	select {
	case w.c <- value:
		ok = true
	case <-w.Dying():
	}
	return
}

func (r readerT[T]) C() <-chan T {
	return r.c
}

func (r *readerT[T]) Read() (value T, ok bool) {
	select {
	case value, ok = <-r.c:
	case <-r.Dying():
	}
	return
}

type pipeOptions[T any] struct {
	parentLifecycle Lifecycle
	childLifecycle  Lifecycle
	childChannel    chan T
}

type PipeOption[T any] func(*pipeOptions[T])

func PipeWithParentLifecycle[T any](parent Lifecycle) PipeOption[T] {
	return func(o *pipeOptions[T]) {
		o.parentLifecycle = parent
	}
}

func PipeWithChildLifecycle[T any](child Lifecycle) PipeOption[T] {
	return func(o *pipeOptions[T]) {
		o.childLifecycle = child
	}
}

func PipeWithChildChannel[T any](c chan T) PipeOption[T] {
	return func(o *pipeOptions[T]) {
		o.childChannel = c
	}
}
