package router

import (
	"hertz-tutorial/demo/biz/stream"

	"github.com/cloudwego/hertz/pkg/route"
)

// registerStream 注册 SSE 流式推送路由，供演示与后续大模型流式输出复用。
func registerStream(g *route.RouterGroup) {
	h := stream.NewHandler()

	s := g.Group("/stream")
	s.GET("/demo", h.Demo) // GET  /api/v1/stream/demo
	s.POST("/chat", h.Chat) // POST /api/v1/stream/chat
}
