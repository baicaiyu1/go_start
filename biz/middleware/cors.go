package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// CORS 允许浏览器从本地 HTML 页面直接调用接口（file:// 或无头调试）。
// 生产环境请把 * 换成具体域名，并谨慎开放 Authorization 等头。
func CORS() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		c.Response.Header.Set("Access-Control-Allow-Origin", "*")
		c.Response.Header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Response.Header.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")
		c.Response.Header.Set("Access-Control-Max-Age", "86400")

		// 预检请求直接返回 204，不再进入业务逻辑
		if string(c.Method()) == "OPTIONS" {
			c.AbortWithStatus(consts.StatusNoContent)
			return
		}
		c.Next(ctx)
	}
}
