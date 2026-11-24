package user

import "encoding/xml"

// CreateUserReq 创建用户的请求结构体（入参 DTO）。
type CreateUserReq struct {
	// 用户名：json 字段名是 name；必填；长度 2~20
	Name string `json:"name" validate:"required,min=2,max=20"`
	// 年龄：用 *int 指针，才能区分“没传”和“传了 0”；必填且范围 0~130
	Age *int `json:"age" validate:"required,gte=0,lte=130"`
	// 邮箱：必填，且必须是合法邮箱
	Email string `json:"email" validate:"required,email"`
	// 角色：缺省时由 handler 填 user；只能是 admin 或 user
	// go-playground 不支持 default 标签，所以这里用 omitempty 允许空值
	Role string `json:"role" validate:"omitempty,oneof=admin user"`
	// 标签：最多 5 个；dive 表示继续校验每个元素；每个元素至少 1 个字符
	Tags []string `json:"tags" validate:"max=5,dive,min=1"`
}

// CreateUserResp 创建用户响应（仅用于 Swagger 文档描述）。
type CreateUserResp struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    CreateUserReq `json:"data"`
	TraceID string        `json:"trace_id,omitempty"`
}

// UserData 用户简要信息。
type UserData struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UserResp 用户简要信息响应（仅用于 Swagger 文档描述）。
type UserResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    UserData `json:"data"`
	TraceID string   `json:"trace_id,omitempty"`
}

// MessageData 通用消息体。
type MessageData struct {
	Message string `json:"message"`
}

// MessageResp 通用消息响应（仅用于 Swagger 文档描述）。
type MessageResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    MessageData `json:"data"`
	TraceID string      `json:"trace_id,omitempty"`
}

// ValidationErrorData 参数错误明细。
type ValidationErrorData struct {
	Fields map[string]string `json:"fields,omitempty"`
	Detail string            `json:"detail,omitempty"`
}

// ValidateErrorResp 绑定/校验失败响应（仅用于 Swagger 文档描述）。
type ValidateErrorResp struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Data    ValidationErrorData `json:"data"`
	TraceID string              `json:"trace_id,omitempty"`
}

// UserXML 用户 XML 结构体。
type UserXML struct {
	XMLName xml.Name `xml:"user"`
	ID      int      `xml:"id"`
	Name    string   `xml:"name"`
}
