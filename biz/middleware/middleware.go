// Package middleware 存放项目自有的 HTTP 中间件。
//
// 约定：
//   - 全局中间件（恢复、链路追踪、访问日志、CORS 等）在 main.go 里通过 h.Use 注册；
//   - 路由组中间件（登录校验、管理员权限、限流等）在 router.go 里通过 group.Use 注册；
//   - 每个中间件一个文件，文件名体现职责，例如 auth.go、cors.go、trace.go、ratelimit.go；
//   - 中间件只依赖 ctx 和 c，不直接写业务逻辑，失败时用 c.Abort() 终止链路。
package middleware

import (
	"context"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/fatih/color"
)

// RequestInfo 打印带颜色的请求日志：方法、路径、状态码、耗时、客户端 IP。
// 在 main.go 中通过 h.Use(middleware.RequestInfo()) 注册。
func RequestInfo() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		c.Next(ctx)

		method := string(c.Method())
		path := c.Path()
		status := c.Response.StatusCode()
		latency := time.Since(start)
		ip := c.ClientIP()

		methodAttr := color.FgGreen
		switch method {
		case "POST":
			methodAttr = color.FgBlue
		case "PUT", "PATCH":
			methodAttr = color.FgYellow
		case "DELETE":
			methodAttr = color.FgRed
		}

		statusAttr := color.FgGreen
		switch {
		case status >= 500:
			statusAttr = color.FgRed
		case status >= 400:
			statusAttr = color.FgYellow
		case status >= 300:
			statusAttr = color.FgCyan
		}

		fmt.Printf(
			"%s %s %s %s %s\n",
			color.New(methodAttr, color.Bold).Sprintf("%-7s", method),
			path,
			color.New(statusAttr, color.Bold).Sprintf("%d", status),
			color.New(color.FgHiBlack).Sprintf("%-12s", latency.Round(time.Microsecond)),
			color.New(color.FgHiBlack).Sprint(ip),
		)
	}
}
