package rx

type Observable[T any] interface {
	Subscribe(subscriber Writer[T])
}

type Writer[T any] interface {
	Lifecycle
	Write(T) bool
}

type Reader[T any] interface {
	Lifecycle
	Read() (T, bool)
}
