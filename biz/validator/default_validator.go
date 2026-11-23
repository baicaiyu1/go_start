package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// NewValidator 创建并配置一个校验器。
// 导出是因为 main.go 需要用它注入 Hertz。
func NewValidator() *validator.Validate {
	// 新建校验器；WithRequiredStructEnabled 让嵌套结构体里的 required 也生效
	vd := validator.New(validator.WithRequiredStructEnabled())
	// 注册字段名函数：让错误信息里显示 json 名，而不是 Go 字段名
	vd.RegisterTagNameFunc(func(fld reflect.StructField) string {
		// 取出 json tag 中逗号前的部分，例如 json:"name,omitempty" 取到 name
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		// 如果没有 json tag，或者 tag 是 -，就退回使用 Go 字段名
		if name == "" || name == "-" {
			return fld.Name
		}
		// 返回 json 字段名
		return name
	})
	// 返回配置好的校验器
	return vd
}