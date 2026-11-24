// Package bootstrap 负责应用启动/关闭阶段的基础设施编排。
//
// 职责边界：
//   - 初始化外部依赖（数据库、Redis、MQ 等），并按逆序关闭；
//   - 管理后台长任务（定时器、消费者、协程）的 context 和退出等待；
//   - 提供启动/下线状态，供健康检查或就绪探针使用。
//
// 约定：OnShutdown 阶段只“通知下线”，真正的资源释放放在 Spin 返回之后。
package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Application 聚合一次进程生命周期内的所有长生命周期资源。
type Application struct {
	resources *Resources
	workers   *WorkerGroup
	ready     atomic.Bool
}

// New 初始化外部依赖和后台任务管理器。
// 任意依赖初始化失败都会回滚已初始化的资源，避免半初始化状态。
func New(ctx context.Context) (*Application, error) {
	app := &Application{
		workers: NewWorkerGroup(ctx),
	}

	resources, err := InitResources(ctx)
	if err != nil {
		app.workers.Cancel()
		return nil, err
	}

	app.resources = resources
	app.ready.Store(true)
	return app, nil
}

// Workers 暴露后台任务管理器，业务在启动阶段注册定时器/消费者。
func (a *Application) Workers() *WorkerGroup {
	return a.workers
}

// Ready 返回当前实例是否可接收流量，可接到 /health 或就绪探针。
func (a *Application) Ready() bool {
	return a.ready.Load()
}

// BeforeShutdown 在 Hertz 的 OnShutdown 阶段执行：
// 只做“准备下线”，不关闭数据库等资源。
func (a *Application) BeforeShutdown(_ context.Context) {
	a.ready.Store(false)
	a.workers.Cancel()
	hlog.Info("[shutdown] 已标记不可用，并通知后台任务停止")
}

// Close 在 h.Spin() 返回后执行：
// 此时 HTTP 已停止接收新请求，在途请求也已处理完，可以安全释放资源。
func (a *Application) Close(ctx context.Context) error {
	var errs []error

	if err := a.workers.Wait(ctx); err != nil {
		errs = append(errs, err)
	}
	if err := a.resources.Close(ctx); err != nil {
		errs = append(errs, err)
	}

	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("应用资源释放失败: %w", err)
	}
	return nil
}
