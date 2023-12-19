package response

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code" example:"20000"` // 代码
	Data    interface{} `json:"data"`                 // 数据集
	Message string      `json:"message"`              // 消息
}

// Error 失败数据处理
func Error(c *gin.Context, err error, response Response) {
	if err != nil && response.Message != "" {
		response.Message = fmt.Sprintf("%s, 错误: %s", response.Message, err.Error())
	} else if err != nil {
		response.Message = err.Error()
	}
	c.AbortWithStatusJSON(http.StatusOK, response)
}

// OK 通常成功数据处理
func OK(c *gin.Context, data interface{}, msg string) {
	var res = Response{
		Code:    Success.Code,
		Data:    data,
		Message: Success.Message,
	}
	if msg != "" {
		res.Message = msg
	}
	if res.Message == "" {
		res.Message = "success"
	}
	c.AbortWithStatusJSON(http.StatusOK, res)
}
