package user

// CreateUserReq 创建用户的请求结构体（入参 DTO）。
type CreateUserReq struct {
	// 用户名：json 字段名是 name；必填；长度 2~20
	Name string `json:"name" validate:"required,min=2,max=20"`
	// 年龄：用 *int 指针，才能区分“没传”和“传了 0”；必填且范围 0~130
	Age *int `json:"age" validate:"required,gte=0,lte=130"`
	// 邮箱：必填，且必须是合法邮箱
	Email string `json:"email" validate:"required,email"`
	// 角色：缺省时由 handler 填 user；只能是 admin 或 user
	// go-playground 不支持 default 标签，所以这里用 omitempty 允许空值
	Role string `json:"role" validate:"omitempty,oneof=admin user"`
	// 标签：最多 5 个；dive 表示继续校验每个元素；每个元素至少 1 个字符
	Tags []string `json:"tags" validate:"max=5,dive,min=1"`
}
