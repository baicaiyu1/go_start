package bootstrap

import (
	"context"
	"errors"
	"fmt"

	"github.com/cloudwego/hertz/pkg/common/hlog"
)

// Resources 持有外部依赖的关闭函数，按注册的逆序释放。
type Resources struct {
	closers []resourceCloser
}

type resourceCloser struct {
	name  string
	close func(context.Context) error
}

// Add 注册一个资源关闭函数。
// 关闭顺序与注册顺序相反：后初始化的先关闭。
func (r *Resources) Add(name string, closeFn func(context.Context) error) {
	if closeFn == nil {
		return
	}
	r.closers = append(r.closers, resourceCloser{name: name, close: closeFn})
}

// Close 逆序关闭所有资源，单个资源失败不影响后续资源关闭。
func (r *Resources) Close(ctx context.Context) error {
	var errs []error

	for i := len(r.closers) - 1; i >= 0; i-- {
		closer := r.closers[i]
		hlog.Infof("[shutdown] 关闭 %s", closer.name)
		if err := closer.close(ctx); err != nil {
			errs = append(errs, fmt.Errorf("关闭 %s 失败: %w", closer.name, err))
		}
	}
	r.closers = nil

	return errors.Join(errs...)
}

// InitResources 按依赖顺序初始化外部资源。
// 初始化失败时回滚已初始化的资源，避免连接泄漏。
func InitResources(ctx context.Context) (*Resources, error) {
	resources := &Resources{}

	initializers := []struct {
		name string
		init func(context.Context) error
	}{
		{name: "database", init: resources.initDB},
		{name: "redis", init: resources.initRedis},
		{name: "mq", init: resources.initMQ},
	}

	for _, item := range initializers {
		if err := item.init(ctx); err != nil {
			if closeErr := resources.Close(ctx); closeErr != nil {
				hlog.Errorf("[startup] 回滚 %s 失败: %v", item.name, closeErr)
			}
			return nil, fmt.Errorf("初始化 %s 失败: %w", item.name, err)
		}
	}

	return resources, nil
}

// initDB 初始化数据库连接池。
func (r *Resources) initDB(_ context.Context) error {
	// 1. 读取数据库配置（DSN、最大连接数、空闲连接数、连接生命周期）
	// 2. 打开连接池：db, err := sql.Open("mysql", dsn)
	// 3. 设置连接池参数：SetMaxOpenConns / SetMaxIdleConns / SetConnMaxLifetime
	// 4. PingContext 检测连通性；失败时先关闭已经创建的连接再返回错误
	// 5. 注册关闭函数：r.Add("database", func(ctx context.Context) error { return db.Close() })
	return nil
}

// initRedis 初始化 Redis 连接。
func (r *Resources) initRedis(_ context.Context) error {
	// 1. 读取 Redis 地址、密码、DB、超时时间
	// 2. 创建客户端：rdb := redis.NewClient(&redis.Options{...})
	// 3. Ping 检测连通性；失败时关闭客户端再返回错误
	// 4. 注册关闭函数：r.Add("redis", func(ctx context.Context) error { return rdb.Close() })
	return nil
}

// initMQ 初始化消息队列生产者/消费者。
func (r *Resources) initMQ(_ context.Context) error {
	// 1. 读取 MQ 地址、Topic、消费组等配置
	// 2. 创建生产者/消费者并建立连接
	// 3. 注册关闭函数：r.Add("mq", func(ctx context.Context) error { return producer.Close() })
	// 4. 消费者如果要长期运行，通过 Application.Workers().Go 启动，不要阻塞这里
	return nil
}
