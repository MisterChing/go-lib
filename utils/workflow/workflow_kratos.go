package workflow

import (
	"context"
	klog "github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"runtime/debug"
	"sync"
)

func KraGoWithRecover(ctx context.Context, logger klog.Logger, fn func(ctx context.Context)) {
	var (
		logH *klog.Helper
	)
	if logger != nil {
		logH = klog.NewHelper(klog.With(logger, "x_module", "utils/GoWithRecover"))
	}
	md, _ := metadata.FromServerContext(ctx)
	ctxCp := metadata.NewServerContext(ctx, md.Clone())
	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				if ctx != nil {
					if logH != nil {
						logH.WithContext(ctx).Errorw(
							"x_msg", "GoWithRecover panic recovered",
							"panic_err", err,
							"stack", string(debug.Stack()),
						)
					}
				} else {
					if logH != nil {
						logH.Errorw(
							"x_msg", "GoWithRecover panic recovered",
							"panic_err", err,
							"stack", string(debug.Stack()),
						)
					}
				}
			}
		}()
		fn(ctx)
	}(ctxCp)
}

func KraGoGroupWait(ctx context.Context, logger klog.Logger, fnArr ...func(ctx context.Context) string) {
	var (
		logH *klog.Helper
	)
	if logger != nil {
		logH = klog.NewHelper(klog.With(logger, "x_module", "utils/GoGroupWait"))
	}
	if len(fnArr) == 0 {
		return
	}
	var wg sync.WaitGroup
	goNum := len(fnArr)
	wg.Add(goNum)
	for _, fn := range fnArr {
		md, _ := metadata.FromServerContext(ctx)
		ctxCp := metadata.NewServerContext(ctx, md.Clone())
		doFn := fn
		go func(ctx context.Context) {
			defer func() {
				if err := recover(); err != nil {
					if ctx != nil {
						if logH != nil {
							logH.WithContext(ctx).Errorw(
								"x_msg", "GoGroupWait panic recovered",
								"panic_err", err,
								"stack", string(debug.Stack()),
							)
						}
					} else {
						if logH != nil {
							logH.Errorw(
								"x_msg", "GoGroupWait panic recovered",
								"panic_err", err,
								"stack", string(debug.Stack()),
							)
						}
					}
				}
			}()
			defer wg.Done()
			doFn(ctx)
		}(ctxCp)
	}
	wg.Wait()
}
