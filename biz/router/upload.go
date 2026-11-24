package router

import (
	"context"

	"hertz-tutorial/demo/biz/upload"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

// registerUpload 注册文件上传路由。
func registerUpload(g *route.RouterGroup) {
	h := upload.NewHandler()

	up := g.Group("/upload")
	up.POST("/one", h.Upload)

	up.POST("/batch", h.BatchUpload)

	// 浏览器跨域预检：CORS 中间件已全局处理，这里保证路由存在
	up.OPTIONS("/", func(ctx context.Context, c *app.RequestContext) {
		c.AbortWithStatus(consts.StatusNoContent)
	})
}
