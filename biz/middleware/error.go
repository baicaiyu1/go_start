package middleware

import (
	"context"

	"hertz-tutorial/demo/biz/response"

	"github.com/cloudwego/hertz/pkg/app"
)

// ErrorHandler 收集 handler 通过 c.Error(err) 抛出的错误，交给 response 统一写回。
// handler 里只需要 c.Error(err) + return，不再各自拼错误 JSON。
func ErrorHandler() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Next(ctx)
		if len(c.Errors) == 0 {
			return
		}
		response.Error(ctx, c, c.Errors.Last().Err)
	}
}
