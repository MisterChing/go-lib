package nerrgroup

import (
	"context"
	"errors"
	"log"
	"runtime/debug"

	"github.com/MisterChing/go-lib/utils/contextutil"
	"golang.org/x/sync/errgroup"
)

var PanicRecoveredError = errors.New("nerrgroup panic recovered")

type Group struct {
	opts options
	ieg  *errgroup.Group
}

func NewGroup(opts ...Option) *Group {
	optsIns := options{}
	for _, o := range opts {
		o(&optsIns)
	}
	ins := &Group{
		opts: optsIns,
		ieg:  &errgroup.Group{},
	}
	if ins.opts.ieg != nil {
		ins.ieg = ins.opts.ieg
	}
	return ins
}

func (ng *Group) wrapWithRecover(fn func() error) func() error {
	return func() (retErr error) {
		defer func() {
			if r := recover(); r != nil {
				retErr = PanicRecoveredError
				if ng.opts.panicFn != nil {
					ng.opts.panicFn(r)
					return
				}
				log.Println(PanicRecoveredError.Error(), string(debug.Stack()))
			}
		}()
		if err := fn(); err != nil {
			retErr = err
		}
		return retErr
	}
}

func (ng *Group) GoWait(fn func() error) {
	ng.ieg.Go(ng.wrapWithRecover(fn))
}

func (ng *Group) GoWaitCtx(ctx context.Context, fn func(ctx context.Context) error) {
	ctxCopy := contextutil.CopyContext(ctx)
	fnCopy := func() error {
		return fn(ctxCopy)
	}
	ng.ieg.Go(ng.wrapWithRecover(fnCopy))
}

func (ng *Group) TryGoWait(fn func() error) {
	ng.ieg.TryGo(ng.wrapWithRecover(fn))
}

func (ng *Group) TryGoWaitCtx(ctx context.Context, fn func(ctx context.Context) error) {
	ctxCopy := contextutil.CopyContext(ctx)
	fnCopy := func() error {
		return fn(ctxCopy)
	}
	ng.ieg.TryGo(ng.wrapWithRecover(fnCopy))
}

func (ng *Group) Wait() error {
	return ng.ieg.Wait()
}

func (ng *Group) GoSafe(fn func() error) {
	go ng.wrapWithRecover(fn)
}

func (ng *Group) GoSafeCtx(ctx context.Context, fn func(ctx context.Context) error) {
	ctxCopy := contextutil.CopyContext(ctx)
	fnCopy := func() error {
		return fn(ctxCopy)
	}
	go ng.wrapWithRecover(fnCopy)
}

func (ng *Group) GoWithoutCancel(ctx context.Context, fn func(ctx context.Context) error) {
	ctxCopy := contextutil.CopyContextWithoutCancel(ctx)
	fnCopy := func() error {
		return fn(ctxCopy)
	}
	go ng.wrapWithRecover(fnCopy)
}

func (ng *Group) SetLimit(n int) {
	ng.ieg.SetLimit(n)
}
