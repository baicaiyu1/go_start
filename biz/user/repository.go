package user

import "context"

// Repository 用户数据访问层。
// 现在用内存 map 演示，以后换成 MySQL/Redis 只改这个文件。
type Repository struct {
	users map[int64]*User
}

// NewRepository 创建仓库（现在是内存版）。
func NewRepository() *Repository {
	return &Repository{
		users: map[int64]*User{
			1: {ID: 1, Name: "八月"},
			2: {ID: 2, Name: "bob"},
		},
	}
}

// FindByID 按 ID 查用户，查不到返回 nil, nil。
func (r *Repository) FindByID(ctx context.Context, id int64) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return r.users[id], nil
}

// Current 返回当前登录用户（演示用，固定 ID=1）。
func (r *Repository) Current(ctx context.Context) (*User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return r.users[1], nil
}
