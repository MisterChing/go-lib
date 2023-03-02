package workflow

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
	var newCtx context.Context
	if ginCtx, ok := kgin.FromGinContext(oldCtx); ok {
		newCtx = kgin.NewGinContext(oldCtx, ginCtx.Copy())
	}
	if md, ok := metadata.FromServerContext(oldCtx); ok {
		newCtx = metadata.NewServerContext(oldCtx, md.Clone())
	}
	return newCtx
}

// CopyContextWithoutCancel 复制kratos context 并且移除cancel
func CopyContextWithoutCancel(oldCtx context.Context) context.Context {
	newCtx := context.Background() //需要重新创建个ctx来避免原ctx cancel
	//0x1 transport context
	if tr, ok := transport.FromServerContext(oldCtx); ok {
		newCtx = transport.NewServerContext(newCtx, tr)
	}
	//0x2 metadata
	if md, ok := metadata.FromServerContext(oldCtx); ok {
		newCtx = metadata.NewServerContext(newCtx, md.Clone())
	}
	//0x3 gin
	if ginCtx, ok := kgin.FromGinContext(oldCtx); ok {
		newCtx = kgin.NewGinContext(newCtx, ginCtx.Copy())
	}
	return newCtx
}
