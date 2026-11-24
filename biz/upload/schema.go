package upload

// UploadData 单文件上传结果。
type UploadData struct {
	Name string `json:"name"`
	Size int64  `json:"size"`
}

// UploadResp 单文件上传响应（仅用于 Swagger 文档描述）。
type UploadResp struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    UploadData `json:"data"`
	TraceID string     `json:"trace_id,omitempty"`
}

// BatchData 批量上传结果。
type BatchData struct {
	Names []string `json:"names"`
}

// BatchResp 批量上传响应（仅用于 Swagger 文档描述）。
type BatchResp struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    BatchData `json:"data"`
	TraceID string    `json:"trace_id,omitempty"`
}

// ErrorResp 上传错误响应（仅用于 Swagger 文档描述）。
type ErrorResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TraceID string `json:"trace_id,omitempty"`
}
