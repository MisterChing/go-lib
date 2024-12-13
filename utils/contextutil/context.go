package contextutil

import (
	"context"

	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
)

func TransferToContext(c *gin.Context) context.Context {
	ctx := metadata.NewServerContext(c.Request.Context(), map[string]string{
		"trace_id": c.GetHeader("trace_id"),
	})
	if _, ok := kgin.FromGinContext(ctx); !ok {
		ctx = kgin.NewGinContext(ctx, c)
		c.Request = c.Request.WithContext(ctx)
	}
	return ctx
}

func CopyContext(oldCtx context.Context) context.Context {
	var (
		newCtx  context.Context
		lastCtx = oldCtx
	)
	//0x1 transport context
	if tr, ok := transport.FromServerContext(oldCtx); ok {
		newCtx = transport.NewServerContext(lastCtx, tr)
		lastCtx = newCtx
	}
	//0x2 metadata
	if md, ok := metadata.FromServerContext(oldCtx); ok {
		newCtx = metadata.NewServerContext(lastCtx, md.Clone())
		lastCtx = newCtx
	}
	//0x3 gin
	if ginCtx, ok := kgin.FromGinContext(oldCtx); ok {
		newCtx = kgin.NewGinContext(lastCtx, ginCtx.Copy())
		lastCtx = newCtx
	}
	if newCtx == nil {
		return oldCtx
	}
	return newCtx
}

// CopyContextWithoutCancel 复制kratos context 并且移除cancel
func CopyContextWithoutCancel(oldCtx context.Context) context.Context {
	var (
		newCtx  = context.Background() //需要重新创建个ctx来避免原ctx cancel
		lastCtx = newCtx
	)
	//0x1 transport context
	if tr, ok := transport.FromServerContext(oldCtx); ok {
		newCtx = transport.NewServerContext(lastCtx, tr)
		lastCtx = newCtx
	}
	//0x2 metadata
	if md, ok := metadata.FromServerContext(oldCtx); ok {
		newCtx = metadata.NewServerContext(lastCtx, md.Clone())
		lastCtx = newCtx
	}
	//0x3 gin
	if ginCtx, ok := kgin.FromGinContext(oldCtx); ok {
		newCtx = kgin.NewGinContext(lastCtx, ginCtx.Copy())
		lastCtx = newCtx
	}
	return newCtx
}
