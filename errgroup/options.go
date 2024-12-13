package errgroup

import (
	syserrgroup "golang.org/x/sync/errgroup"
)

type Option func(*options)

type options struct {
	panicFn func(pac any)
	ieg     *syserrgroup.Group
}

func WithPanicFn(fn func(pac any)) Option {
	return func(o *options) {
		o.panicFn = fn
	}
}

func WithRawGroup(eg *syserrgroup.Group) Option {
	return func(o *options) {
		o.ieg = eg
	}
}
