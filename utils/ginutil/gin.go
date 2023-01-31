package ginutil

import (
	"context"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/metadata"
)

func TransferToContext(c *gin.Context) context.Context {
	ctx := metadata.NewServerContext(c.Request.Context(), map[string]string{
		"xxx": c.GetHeader("xxx"),
	})
	if _, ok := kgin.FromGinContext(ctx); !ok {
		ctx = kgin.NewGinContext(ctx, c)
		c.Request = c.Request.WithContext(ctx)
	}
	return ctx
}
