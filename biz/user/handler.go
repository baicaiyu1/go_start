package user

import (
	"context"
	"errors"
	"strconv"

	"hertz-tutorial/demo/biz/middleware"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Handler 用户模块的 HTTP 层：只做绑定、调 service、写响应。
type Handler struct {
	svc *Service
}

// NewHandler 注入业务层。
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Create POST /api/v1/users
func (h *Handler) Create(ctx context.Context, c *app.RequestContext) {
	var req CreateUserReq
	if err := c.BindAndValidate(&req); err != nil {
		middleware.WriteBindError(c, err)
		return
	}
	// 角色没传时补默认值（go-playground 不支持 default 标签）
	if req.Role == "" {
		req.Role = "user"
	}
	c.JSON(consts.StatusCreated, utils.H{"data": req})
}

// GetByID GET /api/v1/users/:id
func (h *Handler) GetByID(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.JSON(consts.StatusBadRequest, utils.H{"message": "id 不合法"})
		return
	}
	u, err := h.svc.Get(ctx, id)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, utils.H{"id": u.ID, "name": u.Name})
}

// GetCurrent GET /api/v1/users/me
func (h *Handler) GetCurrent(ctx context.Context, c *app.RequestContext) {
	u, err := h.svc.Current(ctx)
	if err != nil {
		h.writeError(c, err)
		return
	}
	c.JSON(consts.StatusOK, utils.H{"id": u.ID, "name": u.Name})
}

// writeError 把业务错误映射成 HTTP 响应。
func (h *Handler) writeError(c *app.RequestContext, err error) {
	if errors.Is(err, ErrUserNotFound) {
		c.JSON(consts.StatusNotFound, utils.H{"message": "用户不存在"})
		return
	}
	c.JSON(consts.StatusInternalServerError, utils.H{"message": "服务器内部错误"})
}
