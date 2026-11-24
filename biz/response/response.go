// Package response 统一 HTTP 响应结构与错误映射。
//
// 约定：
//   - 成功统一走 OK / Created，响应体为 APIResponse；
//   - 绑定错误在 handler 里调用 BindError；
//   - 业务错误通过 c.Error(err) 交给 middleware.ErrorHandler，再调用 Error 写回；
//   - trace_id 由 middleware.Trace 注入响应头，这里负责回填到响应体。
package response

import (
	"hertz-tutorial/demo/biz/errs"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// TraceIDHeader 链路追踪请求/响应头。
const TraceIDHeader = "X-Trace-ID"

// APIResponse 全局统一响应结构。
// code=0 表示成功，非 0 为业务码；HTTP 状态码继续表达传输层语义。
type APIResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	TraceID string `json:"trace_id,omitempty"`
}

// TraceID 读取本次请求的链路 ID，未启用 Trace 中间件时返回空。
func TraceID(c *app.RequestContext) string {
	return string(c.Response.Header.Peek(TraceIDHeader))
}

// OK 返回 200 成功响应。
func OK(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, APIResponse{
		Code:    errs.CodeOK,
		Message: "ok",
		Data:    data,
		TraceID: TraceID(c),
	})
}

// Created 返回 201 创建成功响应。
func Created(c *app.RequestContext, data any) {
	c.JSON(consts.StatusCreated, APIResponse{
		Code:    errs.CodeOK,
		Message: "ok",
		Data:    data,
		TraceID: TraceID(c),
	})
}

// Fail 写入失败响应并终止后续 handler。
func Fail(c *app.RequestContext, httpStatus, code int, message string) {
	c.AbortWithStatusJSON(httpStatus, APIResponse{
		Code:    code,
		Message: message,
		TraceID: TraceID(c),
	})
}

// failWithData 写入带 data 的失败响应，供参数校验等需要返回明细的场景使用。
func failWithData(c *app.RequestContext, httpStatus, code int, message string, data any) {
	c.AbortWithStatusJSON(httpStatus, APIResponse{
		Code:    code,
		Message: message,
		Data:    data,
		TraceID: TraceID(c),
	})
}
