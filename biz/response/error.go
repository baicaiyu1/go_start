package response

import (
	"context"
	"errors"
	"strings"

	"hertz-tutorial/demo/biz/errs"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/go-playground/validator/v10"
)

// validationData 参数错误明细，放在统一响应的 data 字段里。
type validationData struct {
	Fields map[string]string `json:"fields,omitempty"`
	Detail string            `json:"detail,omitempty"`
}

// BindError 处理 c.BindAndValidate 的返回值：
// 校验失败返回 422 + 字段规则；JSON 格式/类型错误返回 400 + 原始信息。
func BindError(c *app.RequestContext, err error) {
	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		fields := make(map[string]string, len(ve))
		for _, fe := range ve {
			rule := fe.Tag()
			if param := fe.Param(); param != "" {
				rule += "=" + param
			}
			// Namespace 形如 CreateUserReq.Address.City，去掉结构体名只留字段路径
			name := fe.Field()
			if ns := fe.Namespace(); strings.Contains(ns, ".") {
				name = ns[strings.Index(ns, ".")+1:]
			}
			fields[name] = rule
		}
		failWithData(c, consts.StatusUnprocessableEntity, errs.CodeValidationFailed, "参数校验失败",
			validationData{Fields: fields})
		return
	}

	failWithData(c, consts.StatusBadRequest, errs.CodeBadRequest, "参数格式错误",
		validationData{Detail: err.Error()})
}

// Error 统一处理业务错误：BizError 按其状态码/业务码输出，
// 未识别的错误记录日志并返回 500，避免把内部细节暴露给调用方。
func Error(ctx context.Context, c *app.RequestContext, err error) {
	var be *errs.BizError
	if errors.As(err, &be) {
		Fail(c, be.HTTPStatus, be.Code, be.Message)
		return
	}

	hlog.CtxErrorf(ctx, "unhandled error: %v", err)
	Fail(c, consts.StatusInternalServerError, errs.CodeInternal, "服务器内部错误")
}
