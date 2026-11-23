package user

import (
	"context"
	"errors"
	"fmt"
)

// ErrUserNotFound 业务错误，由 handler 映射成 HTTP 404。
var ErrUserNotFound = errors.New("user not found")

// Service 用户业务层：写业务规则，不关心 HTTP，也不关心 SQL。
type Service struct {
	repo *Repository
}

// NewService 注入仓库。
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Get 业务规则：ID 必须合法，查不到统一返回 ErrUserNotFound。
func (s *Service) Get(ctx context.Context, id int64) (*User, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid user id: %d", id)
	}
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}

// Current 获取当前登录用户。
func (s *Service) Current(ctx context.Context) (*User, error) {
	u, err := s.repo.Current(ctx)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ErrUserNotFound
	}
	return u, nil
}
