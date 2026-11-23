package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetFile GET /api/v1/files/*filepath
func GetFile(ctx context.Context, c *app.RequestContext) {
	c.JSON(consts.StatusOK, utils.H{"filepath": c.Param("filepath")})
}
