package workflow

import (
	"context"
	"github.com/MisterChing/go-lib/utils/debugutil"
	"github.com/go-kratos/kratos/v2/metadata"
	"testing"
)

func TestCopyContext(t *testing.T) {
	var (
		ctx2 context.Context
	)
	ctx := context.Background()
	ctx = metadata.NewServerContext(context.Background(), nil)
	ctcCp := CopyContext(ctx)

	debugutil.DebugPrintV2("111", ctx, ctcCp, ctx2)

}
