package repository

import (
	"context"

	"hertz-tutorial/demo/biz/model"
)

// users 演示用的内存数据，后续换成 MySQL / Redis 时只改这一层。
var users = map[int64]*model.User{
	1: {ID: 1, Name: "八月"},
	2: {ID: 2, Name: "bob"},
}

// currentUserID 演示用的当前登录用户。
const currentUserID int64 = 1

// FindUserByID 按 ID 查用户，查不到返回 nil, nil。
func FindUserByID(ctx context.Context, id int64) (*model.User, error) {
	return users[id], nil
}

// CurrentUser 返回当前登录用户。
func CurrentUser(ctx context.Context) (*model.User, error) {
	return users[currentUserID], nil
}
