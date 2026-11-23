package router

import (
	"hertz-tutorial/demo/biz/system"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/route"
)

// registerSystem 注册基础设施路由，不挂版本前缀和业务中间件。
func registerSystem(r *server.Hertz) {
	r.GET("/ping", system.Ping)
	r.GET("/health", system.Health)
}

// registerFile 注册文件路由（跟着 v1 版本组走）。
func registerFile(g *route.RouterGroup) {
	g.GET("/files/*filepath", system.GetFile)
}
