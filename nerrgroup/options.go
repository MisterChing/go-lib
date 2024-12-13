package nerrgroup

import (
	"golang.org/x/sync/errgroup"
)

type Option func(*options)

type options struct {
	panicFn func(pac any)
	ieg     *errgroup.Group
}

func WithPanicFn(fn func(pac any)) Option {
	return func(o *options) {
		o.panicFn = fn
	}
}

func WithRawGroup(eg *errgroup.Group) Option {
	return func(o *options) {
		o.ieg = eg
	}
}
