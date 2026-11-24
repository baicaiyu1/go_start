// Package system 放与具体业务无关的基础设施接口：健康检查、静态文件等。
package system

import (
	"context"
	"sync/atomic"

	"hertz-tutorial/demo/biz/errs"
	"hertz-tutorial/demo/biz/response"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// readinessCheck 由 main 注入，用于关闭阶段让 /health 返回 503。
var readinessCheck atomic.Value // func() bool

func init() {
	readinessCheck.Store(func() bool { return true })
}

// SetReadinessCheck 注入就绪检查函数，传 nil 会忽略。
func SetReadinessCheck(fn func() bool) {
	if fn != nil {
		readinessCheck.Store(fn)
	}
}

// Ping GET /ping
//
// @Summary 连通性探测
// @Tags 基础设施
// @Produce json
// @Success 200 {object} MessageResp
// @Router /ping [get]
func Ping(ctx context.Context, c *app.RequestContext) {
	response.OK(c, utils.H{"message": "pong"})
}

// Health GET /health
//
// @Summary 健康检查
// @Tags 基础设施
// @Produce json
// @Success 200 {object} MessageResp
// @Failure 503 {object} MessageResp
// @Router /health [get]
func Health(ctx context.Context, c *app.RequestContext) {
	if fn, ok := readinessCheck.Load().(func() bool); ok && !fn() {
		response.Fail(c, consts.StatusServiceUnavailable, errs.CodeUnavailable, "服务正在下线")
		return
	}
	response.OK(c, utils.H{"message": "healthy"})
}

// GetFile GET /api/v1/files/*filepath
//
// @Summary 回显文件路径
// @Tags 基础设施
// @Produce json
// @Param filepath path string true "文件路径"
// @Success 200 {object} FileResp
// @Router /api/v1/files/{filepath} [get]
func GetFile(ctx context.Context, c *app.RequestContext) {
	response.OK(c, utils.H{"filepath": c.Param("filepath")})
}
