package pagination

import (
	"strings"

	"github.com/gin-gonic/gin"
)

/*
  @Author : lanyulei
*/

func RequestParams(c *gin.Context) map[string]interface{} {
	params := make(map[string]interface{}, 10)
	queryList := strings.Split(c.Request.URL.RawQuery, "&")
	if len(queryList) > 0 {
		for _, value := range queryList {
			values := strings.Split(value, "=")

			if values[0] == "page" || values[0] == "size" || values[0] == "sort" || values[0] == "not_page" {
				continue
			}
			params[values[0]] = values[1]
		}
	}
	return params
}
