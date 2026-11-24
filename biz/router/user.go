package router

import (
	"context"

	"hertz-tutorial/demo/biz/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
)

// registerUser 注册用户领域路由。
// 依赖在这里装配：repository -> service -> handler。
func registerUser(g *route.RouterGroup) {
	repo := user.NewRepository()
	svc := user.NewService(repo)
	h := user.NewHandler(svc)

	users := g.Group("/users")
	users.POST("", h.Create)       // POST   /api/v1/users
	users.GET("/me", h.GetCurrent) // GET    /api/v1/users/me（静态路由写在 :id 前面）
	// users.GET("/:id", h.GetByID)   // GET    /api/v1/users/:id
	users.GET("/index", h.IndexHTML)
	users.GET("/old", h.RedirectIndex)
	users.GET("/:id/xml", h.GetByIDXML)
	users.GET("/:id", h.GetByIDJSON)

	// 浏览器跨域预检：CORS 中间件已全局处理，这里保证路由存在
	users.OPTIONS("", func(ctx context.Context, c *app.RequestContext) {
		c.AbortWithStatus(consts.StatusNoContent)
	})
}
