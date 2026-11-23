package middleware

import (
	"errors"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/go-playground/validator/v10"
)

// WriteBindError 把绑定/校验错误写成统一的 JSON 响应。
// 导出是因为 handler 包需要调用它。
func WriteBindError(c *app.RequestContext, err error) {
	// 声明一个校验错误集合，用来接收 validator 返回的错误
	var ve validator.ValidationErrors
	// errors.As 判断 err 是否属于校验错误；是的话把值放进 ve
	if errors.As(err, &ve) {
		// 创建 map，key 是字段名，value 是失败规则
		fields := make(map[string]string, len(ve))
		// 遍历每一个字段错误
		for _, fe := range ve {
			// 取出失败规则名，例如 min、email、gte
			rule := fe.Tag()
			// 如果规则带参数（例如 min=2），把参数拼上
			if param := fe.Param(); param != "" {
				rule += "=" + param
			}
			// 取出字段名（前面注册了 TagNameFunc，所以这里是 json 名）
			name := fe.Field()
			// 嵌套字段用 Namespace 保留层级，例如 address.city
			if ns := fe.Namespace(); strings.Contains(ns, ".") {
				// 去掉最前面的结构体名，只保留 address.city
				name = ns[strings.Index(ns, ".")+1:]
			}
			// 记录“字段名 -> 失败规则”
			fields[name] = rule
		}
		// 返回 422，并带上 code/message/fields
		c.AbortWithStatusJSON(consts.StatusUnprocessableEntity, utils.H{
			"code":    422,           // 业务码：参数校验失败
			"message": "参数校验失败", // 给用户看的提示
			"fields":  fields,        // 每个字段的具体失败规则
		})
		// 结束当前 handler
		return
	}
	// 不是校验错误，就当作格式/类型错误，返回 400
	c.AbortWithStatusJSON(consts.StatusBadRequest, utils.H{
		"code":    400,         // 业务码：请求格式错误
		"message": "参数格式错误", // 给用户看的提示
		"detail":  err.Error(), // 原始错误信息，方便排查
	})
}