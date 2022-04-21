package rx

import (
	"context"

	"gopkg.in/tomb.v2"
)

type Lifecycle interface {
	Kill(error)
	Dead() <-chan struct{}
	Dying() <-chan struct{}
	Wait() error
	Go(func() error)
	Err() error
	Alive() bool
	Context(context.Context) context.Context
}

func NewLifecycle() Lifecycle {
	return &tomb.Tomb{}
}

func Bind(parent, child Lifecycle, options ...BindOption) {
	if parent == nil || child == nil {
		return
	}
	opts := buildOptions(options)
	parent.Go(func() (err error) {
		select {
		case <-child.Dying():
			if opts.discardChildError {
				child.Wait()
			} else {
				err = child.Wait()
			}
		case <-parent.Dying():
			child.Kill(parent.Err())
		}
		return
	})
}

func BindWithDiscardChildError() BindOption {
	return func(o *bindOptions) {
		o.discardChildError = true
	}
}

type BindOption func(*bindOptions)

type bindOptions struct {
	discardChildError bool
}
