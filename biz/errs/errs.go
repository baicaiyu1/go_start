// Package errs 定义业务错误类型与错误码。
//
// 约定：
//   - BizError 同时携带 HTTP 状态码和业务码，service 只返回错误，不拼 HTTP 响应；
//   - 错误由 middleware.ErrorHandler 统一映射成 JSON；
//   - 需要保留底层错误时用 Wrap，errors.Is / errors.As 仍可沿 cause 向上匹配。
package errs

// BizError 业务错误：HTTPStatus 决定响应状态码，Code 给前端做分支，Message 给用户看。
type BizError struct {
	HTTPStatus int
	Code       int
	Message    string
	cause      error
}

// Error 实现 error 接口，只返回用户可读信息。
func (e *BizError) Error() string { return e.Message }

// Unwrap 支持 errors.Is / errors.As 继续匹配底层错误。
func (e *BizError) Unwrap() error { return e.cause }

// New 创建一个不带底层错误的业务错误。
func New(httpStatus, code int, message string) *BizError {
	return &BizError{HTTPStatus: httpStatus, Code: code, Message: message}
}

// Wrap 包装底层错误，保留原始信息用于日志排查。
func Wrap(cause error, httpStatus, code int, message string) *BizError {
	return &BizError{HTTPStatus: httpStatus, Code: code, Message: message, cause: cause}
}
