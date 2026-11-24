package bootstrap

import (
	"context"
	"fmt"
	"sync"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// WorkerGroup 管理后台长任务，统一接收取消信号并等待退出。
type WorkerGroup struct {
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewWorkerGroup 创建后台任务管理器，ctx 作为所有任务的父 context。
func NewWorkerGroup(parent context.Context) *WorkerGroup {
	ctx, cancel := context.WithCancel(parent)
	return &WorkerGroup{ctx: ctx, cancel: cancel}
}

// Context 返回后台任务共用的 context，任务内部应监听 Done()。
func (g *WorkerGroup) Context() context.Context {
	return g.ctx
}

// Go 启动一个后台任务；任务退出或 panic 都不会影响主进程。
func (g *WorkerGroup) Go(name string, fn func(context.Context)) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		defer func() {
			if recovered := recover(); recovered != nil {
				hlog.Errorf("[worker] %s panic: %v", name, recovered)
			}
		}()
		fn(g.ctx)
	}()
}

// Cancel 通知所有后台任务停止接新任务，任务应尽快收尾退出。
func (g *WorkerGroup) Cancel() {
	g.cancel()
}

// Wait 等待所有后台任务退出；超过 ctx 截止时间则返回超时错误。
func (g *WorkerGroup) Wait(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		g.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("等待后台任务退出超时: %w", ctx.Err())
	}
}
