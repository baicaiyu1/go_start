package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/fatih/color"
)

// AccessLog 打印带颜色的访问日志：方法、路径、状态码、耗时、客户端 IP。
// 必须注册在 ErrorHandler 之前，才能记录到错误中间件写回后的真实状态码。
func AccessLog() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Path())
		// 健康检查会被高频轮询，Swagger 是开发期工具，都不记录，避免日志刷屏
		if path == "/ping" || path == "/health" || strings.HasPrefix(path, "/swagger") {
			c.Next(ctx)
			return
		}

		start := time.Now()
		c.Next(ctx)

		method := string(c.Method())
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
