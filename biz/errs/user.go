package errs

import "net/http"

// 用户域错误码：404 下的第 1 个。
const CodeUserNotFound = 40401

// ErrUserNotFound 用户不存在，service 直接返回，由中间件映射成 404。
var ErrUserNotFound = New(http.StatusNotFound, CodeUserNotFound, "用户不存在")
