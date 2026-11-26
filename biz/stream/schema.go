package stream

// ChatReq 流式对话请求：后续接大模型时，把这里换成真实的 messages 结构即可。
type ChatReq struct {
	Prompt string `json:"prompt" validate:"required,min=1,max=2000"`
}

// ErrorResp 流式接口的错误响应（仅用于 Swagger 文档描述）。
type ErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
}
