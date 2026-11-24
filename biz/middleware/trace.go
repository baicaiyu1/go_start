package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"strings"
	"time"

	"hertz-tutorial/demo/biz/response"

	"github.com/cloudwego/hertz/pkg/app"
)

// Trace 为每个请求生成链路 ID：优先复用上游传入的 X-Trace-ID，
// 然后写回响应头，response 包会自动带进 JSON。
func Trace() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		traceID := strings.TrimSpace(string(c.Request.Header.Peek(response.TraceIDHeader)))
		if traceID == "" || len(traceID) > 64 {
			traceID = newTraceID()
		}
		c.Response.Header.Set(response.TraceIDHeader, traceID)
		c.Next(ctx)
	}
}

func newTraceID() string {
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		// 随机数不可用属于极端情况，退化成时间戳，保证链路 ID 不中断
		return "t" + time.Now().Format("20060102150405.000000000")
	}
	return hex.EncodeToString(buf)
}
