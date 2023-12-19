package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	v1 "github.com/lanyulei/comet/common/router/v1"
	"github.com/lanyulei/comet/pkg/logger"
	"github.com/lanyulei/comet/pkg/tools/response"
	"net/http"
)

// Setup 加载路由
func Setup(g *gin.Engine) {
	// 使用zap接收gin框架默认的日志并配置
	g.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 404
	g.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    http.StatusNotFound,
			"message": "404 page not found",
		})
	})

	// pprof router
	pprof.Register(g)

	// 健康检查接口
	g.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code":    response.Success.Code,
			"message": response.Success.Message,
		})
	})

	// 路由版本
	v1.RegisterRouter(g.Group(ApiV1Version))
}
