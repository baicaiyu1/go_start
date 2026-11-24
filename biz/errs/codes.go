package errs

// 通用错误码：前 3 位等于 HTTP 状态码，后 2 位是域内序号。
// 例如 42200 表示 422 下的通用参数校验错误，40401 表示 404 下的第 1 个业务错误。
const (
	CodeOK = 0

	CodeBadRequest       = 40000
	CodeUnauthorized     = 40100
	CodeForbidden        = 40300
	CodeNotFound         = 40400
	CodeMethodNotAllowed = 40500
	CodeConflict         = 40900
	CodeValidationFailed = 42200
	CodeInternal         = 50000
	CodeUnavailable      = 50300
)
