package auth

import (
	"comet/pkg/respstatus"
	"github.com/lanyulei/toolkit/jwtauth"
	"github.com/lanyulei/toolkit/response"
	"github.com/spf13/viper"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			response.Error(c, nil, respstatus.AuthorizationNullError)
			c.Abort()
			return
		}

		// 按空格分割
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Error(c, nil, respstatus.AuthorizationFormatError)
			c.Abort()
			return
		}

		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwtauth.ParseToken(parts[1], viper.GetString("jwt.secret"))
		if err != nil {
			response.Error(c, err, respstatus.InvalidTokenError)
			c.Abort()
			return
		}

		// 将当前请求的username信息保存到请求的上下文c上
		c.Set("username", mc.Username)
		c.Set("userId", mc.UserId)
		c.Set("isAdmin", mc.IsAdmin)
		c.Next() // 后续的处理函数可以用过c.Get("username")来获取当前请求的用户信息
	}
}
