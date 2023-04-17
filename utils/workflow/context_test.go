package workflow

import (
	"context"
	"github.com/MisterChing/go-lib/utils/debugutil"
	"github.com/gin-gonic/gin"
	kgin "github.com/go-kratos/gin"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"testing"
)

func TestCopyContext(t *testing.T) {
	var (
		ctx2 context.Context
	)
	ctx := context.Background()
	ctx = transport.NewServerContext(ctx, &http.Transport{})
	ctx = metadata.NewServerContext(ctx, metadata.New(map[string]string{
		"aa": "aa",
	}))
	ctx = kgin.NewGinContext(ctx, &gin.Context{})
	md, _ := metadata.FromServerContext(ctx)
	debugutil.DebugPrintV2("before", ctx, ctx2, md)

	ctxCp := CopyContext(ctx)
	ctxCpWithNoCancel := CopyContextWithoutCancel(ctx)
	md1, _ := metadata.FromServerContext(ctx)
	md2, _ := metadata.FromServerContext(ctxCp)
	md3, _ := metadata.FromServerContext(ctxCpWithNoCancel)
	debugutil.DebugPrintV2("after", ctx, ctxCp, ctxCpWithNoCancel, ctx2, md1, md2, md3)

}
