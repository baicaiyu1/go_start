// Package system 放与具体业务无关的基础设施接口：健康检查、静态文件等。
package system

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Ping GET /ping
func Ping(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"message": "pong"})
}

// Health GET /health
func Health(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"message": "healthy"})
}

// GetFile GET /api/v1/files/*filepath
func GetFile(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"filepath": c.Param("filepath")})
}
