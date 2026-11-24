package middleware

import (
	"context"
	"net/http"

	"hertz-tutorial/demo/biz/errs"
	"hertz-tutorial/demo/biz/response"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/middlewares/server/recovery"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Recovery 捕获 panic：记录堆栈，并返回统一的 500 JSON，避免暴露裸文本或导致进程退出。
// 必须作为第一个全局中间件注册，才能兜住后面所有中间件和 handler 的 panic。
func Recovery() app.HandlerFunc {
	return recovery.Recovery(recovery.WithRecoveryHandler(
		func(ctx context.Context, c *app.RequestContext, err any, stack []byte) {
			hlog.CtxErrorf(ctx, "panic recovered: %v\n%s", err, stack)
			response.Fail(c, http.StatusInternalServerError, errs.CodeInternal, "服务器内部错误")
		},
	))
}
