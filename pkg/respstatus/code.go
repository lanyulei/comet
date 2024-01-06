package respstatus

import "github.com/lanyulei/toolkit/response"

var (
	AuthorizationNullError = response.Response{Code: 30001, Message: "请求头中 Authorization 为空"}
	InvalidTokenError      = response.Response{Code: 30002, Message: "Token 无效"}
	UnknownError           = response.Response{Code: 30005, Message: "未知错误"}

	AuthorizationFormatError = response.Response{Code: 40135, Message: "验证格式错误"}
)
