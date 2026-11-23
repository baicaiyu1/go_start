package handler

import (
	"context"
	"errors"
	"strconv"

	"hertz-tutorial/demo/biz/service"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// GetCurrentUser GET /api/v1/users/me
func GetCurrentUser(ctx context.Context, c *app.RequestContext) {
	user, err := service.GetCurrentUser(ctx)
	if err != nil {
		writeUserError(c, err)
		return
	}
	c.JSON(consts.StatusOK, utils.H{"id": user.ID, "name": user.Name})
}

// GetUserByID GET /api/v1/users/:id
func GetUserByID(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "id 不合法"})
		return
	}
	user, err := service.GetUser(ctx, id)
	if err != nil {
		writeUserError(c, err)
		return
	}
	c.JSON(consts.StatusOK, utils.H{"id": user.ID, "name": user.Name})
}

// writeUserError 把业务错误统一映射成 HTTP 响应。
func writeUserError(c *app.RequestContext, err error) {
	if errors.Is(err, service.ErrUserNotFound) {
		c.JSON(consts.StatusNotFound, utils.H{"message": "用户不存在"})
		return
	}
	c.JSON(consts.StatusInternalServerError, utils.H{"message": "服务器内部错误"})
}
