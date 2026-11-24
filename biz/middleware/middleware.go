// Package middleware 存放项目自有的 HTTP 中间件。
//
// 文件职责：
//   - recovery.go  panic 恢复，必须作为第一个全局中间件注册；
//   - trace.go     生成/透传 X-Trace-ID，统一响应里的 trace_id 来自这里；
//   - accesslog.go 访问日志，注册在错误收敛之前，跳过健康检查和 Swagger；
//   - error.go     c.Error 错误收敛；
//   - cors.go      跨域。
//
// 注册顺序（main.go）：Recovery -> Trace -> AccessLog -> ErrorHandler -> CORS。
// 约定：中间件只依赖 ctx 和 c，不写业务逻辑；失败时 Abort 终止链路。
package middleware
