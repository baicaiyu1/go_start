package bootstrap

import (
	"context"
	"net"
	"time"

	limits "hertz-tutorial/demo/biz/consts"
	"hertz-tutorial/demo/biz/middleware"
	"hertz-tutorial/demo/biz/validator"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/network"
	"github.com/cloudwego/hertz/pkg/network/standard"
	"github.com/cloudwego/hertz/pkg/protocol"
)

// NewHTTP 装配 HTTP 服务：连接参数、参数校验、连接级 Hook、启动/退出 Hook 和全局中间件。
// 路由注册不放进这里，main 仍走生成代码的 register 入口，避免出现两套注册方式。
func NewHTTP(app *Application) *server.Hertz {
	vd := validator.NewValidator()

	h := server.New(
		server.WithHostPorts("0.0.0.0:8080"),
		server.WithReadTimeout(10*time.Second),
		server.WithWriteTimeout(10*time.Second),
		server.WithIdleTimeout(60*time.Second),
		server.WithMaxRequestBodySize(limits.MaxRequestBodySize),
		server.WithExitWaitTime(10*time.Second),
		server.WithHandleMethodNotAllowed(true),
		server.WithTransport(standard.NewTransporter),
		// OnAccept：连接刚建立、还不能读写数据，每个连接一次
		server.WithOnAccept(func(conn net.Conn) context.Context {
			hlog.Infof("[hook] OnAccept 连接建立 conn=%s", conn.RemoteAddr())
			// TODO: 连接级限流、IP 黑名单、连接打标
			return context.Background()
		}),
		// OnConnect：连接可读写后（标准库下是 TLS 握手后），每个连接一次
		server.WithOnConnect(func(ctx context.Context, conn network.Conn) context.Context {
			hlog.Infof("[hook] OnConnect 连接可读 conn=%s", conn.RemoteAddr())
			// TODO: 协议嗅探（Peek 首字节判断 TLS/HTTP）、连接级鉴权、读取首包
			return ctx
		}),
		server.WithCustomValidatorFunc(func(_ *protocol.Request, req any) error {
			// 对传入的结构体执行校验
			return vd.Struct(req)
		}),
	)

	// OnRun：监听之前顺序执行，返回 error 会终止启动
	h.Engine.OnRun = append(h.Engine.OnRun, func(ctx context.Context) error {
		hlog.Info("[hook] OnRun 开始：服务启动前")
		// TODO: 校验关键配置项，缺失或非法直接返回 error 拒绝启动
		// TODO: 预热热点缓存、初始化下游连接池
		// TODO: 注册到注册中心（先注册成功，再开始对外提供服务）
		hlog.Info("[hook] OnRun 完成")
		return nil
	})

	// OnShutdown 是并发执行的，多个 Hook 之间没有先后顺序；
	// 所以“置不可用 -> 摘流量 -> 停消费 -> flush”必须放在同一个 Hook 内顺序写
	h.Engine.OnShutdown = append(h.Engine.OnShutdown, func(ctx context.Context) {
		// 只做“准备下线”：标记不可用、通知后台任务停止
		// 数据库/Redis/MQ 的关闭放在 h.Spin() 返回之后
		app.BeforeShutdown(ctx)

		hlog.Info("[hook] OnShutdown 收尾开始")
		// TODO: 从注册中心摘除，让调用方尽快停止打流量
		// TODO: 停止 Kafka/MQ 消费者并提交位点，避免重复消费
		// TODO: flush 日志、指标、trace，避免最后一批数据丢失
		hlog.Info("[hook] OnShutdown 收尾完成")
	})

	// 全局中间件顺序：panic 恢复最外层，Trace 尽早注入，访问日志在错误收敛之前
	h.Use(middleware.Recovery())
	h.Use(middleware.Trace())
	h.Use(middleware.AccessLog())
	h.Use(middleware.ErrorHandler())
	h.Use(middleware.CORS())

	return h
}
