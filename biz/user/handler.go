package user

import (
	"context"
	"strconv"

	"hertz-tutorial/demo/biz/errs"
	"hertz-tutorial/demo/biz/response"

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
//
// @Summary 创建用户
// @Description 校验请求体后回显用户信息；role 缺省为 user
// @Tags 用户
// @Accept json
// @Produce json
// @Param req body CreateUserReq true "用户信息"
// @Success 201 {object} CreateUserResp
// @Failure 400 {object} ValidateErrorResp
// @Failure 422 {object} ValidateErrorResp
// @Router /api/v1/users [post]
func (h *Handler) Create(ctx context.Context, c *app.RequestContext) {
	var req CreateUserReq
	if err := c.BindAndValidate(&req); err != nil {
		response.BindError(c, err)
		return
	}
	// 角色没传时补默认值（go-playground 不支持 default 标签）
	if req.Role == "" {
		req.Role = "user"
	}
	response.Created(c, req)
}

// GetByID GET /api/v1/users/:id
//
// @Summary 按 ID 查询用户
// @Tags 用户
// @Produce json
// @Param id path int true "用户 ID"
// @Success 200 {object} UserResp
// @Failure 400 {object} MessageResp
// @Failure 404 {object} MessageResp
// @Router /api/v1/users/{id} [get]
func (h *Handler) GetByID(ctx context.Context, c *app.RequestContext) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		c.Error(errs.New(consts.StatusBadRequest, errs.CodeBadRequest, "id 不合法"))
		return
	}
	u, err := h.svc.Get(ctx, id)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, utils.H{"id": u.ID, "name": u.Name})
}

// GetCurrent GET /api/v1/users/me
//
// @Summary 查询当前用户
// @Tags 用户
// @Produce json
// @Success 200 {object} UserResp
// @Failure 404 {object} MessageResp
// @Router /api/v1/users/me [get]
func (h *Handler) GetCurrent(ctx context.Context, c *app.RequestContext) {
	u, err := h.svc.Current(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	response.OK(c, utils.H{"id": u.ID, "name": u.Name})
}

func (h *Handler) GetByIDJSON(ctx context.Context, c *app.RequestContext) {
	response.OK(c, utils.H{"id": c.Param("id"), "msg": "success"})
}

func (h *Handler) GetByIDXML(ctx context.Context, c *app.RequestContext) {
	c.XML(consts.StatusOK, UserXML{ID: 1, Name: "tom"})
}

func (h *Handler) IndexHTML(ctx context.Context, c *app.RequestContext) {
	c.HTML(consts.StatusOK, "index.tmpl", utils.H{"title": "Hertz", "name": "tom"})
}

func (h *Handler) RedirectIndex(ctx context.Context, c *app.RequestContext) {
	c.Redirect(consts.StatusFound, []byte("/api/v1/users/index"))
}
