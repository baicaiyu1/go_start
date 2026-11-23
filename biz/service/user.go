package service

import (
	"context"
	"errors"
	"fmt"

	"hertz-tutorial/demo/biz/model"
	"hertz-tutorial/demo/biz/repository"
)

// ErrUserNotFound 业务错误，由 handler 映射成 HTTP 状态码。
var ErrUserNotFound = errors.New("user not found")

// GetUser 业务规则：ID 必须合法，查不到统一返回 ErrUserNotFound。
func GetUser(ctx context.Context, id int64) (*model.User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", id)
	}
	user, err := repository.FindUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// GetCurrentUser 获取当前登录用户。
func GetCurrentUser(ctx context.Context) (*model.User, error) {
	user, err := repository.CurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}
