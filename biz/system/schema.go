package system

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

// FileData 文件路径回显。
type FileData struct {
	Filepath string `json:"filepath"`
}

// FileResp 文件路径回显响应（仅用于 Swagger 文档描述）。
type FileResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Data    FileData `json:"data"`
	TraceID string   `json:"trace_id,omitempty"`
}
