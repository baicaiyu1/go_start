package user

// User 用户实体，只描述数据，不含业务行为。
// 以后加 MySQL 字段（Email、Age、CreatedAt）也加在这里。
type User struct {
	ID   int64
	Name string
}
