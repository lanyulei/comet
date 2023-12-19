package pagination

/*
 @Author : lanyulei
*/

import (
	"fmt"
	"github.com/lanyulei/comet/pkg/logger"
	"math"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Param struct {
	C        *gin.Context
	DB       *gorm.DB
	ShowSQL  bool
	Distinct bool
}

type Paginator struct {
	Total     int64       `json:"total"`
	TotalPage int         `json:"total_page"`
	List      interface{} `json:"list"`
	Size      int         `json:"size"`
	Page      int         `json:"page"`
}

type ListRequest struct {
	NotPage bool `json:"not_page" form:"not_page"`
	Page    int  `json:"page" form:"page"`
	Size    int  `json:"size" form:"size"`
	Sort    int  `json:"sort" form:"sort"`
}

// Paging 分页
func Paging(p *Param, result interface{}, args ...interface{}) (*Paginator, error) {
	var (
		err       error
		param     ListRequest
		paginator Paginator
		count     int64
		offset    int
		tableName string
	)

	if err = p.C.Bind(&param); err != nil {
		logger.Errorf("参数绑定失败，错误：%v", err)
		return nil, err
	}

	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}

	if len(args) > 1 {
		tableName = fmt.Sprintf("\"%s\".", args[1].(string))
	}

	if len(args) > 0 && args[0] != nil {
		for paramType, paramsValue := range args[0].(map[string]map[string]interface{}) {
			if paramType == "like" {
				for key, value := range paramsValue {
					if key != "" {
						db = db.Where(fmt.Sprintf("%v%v like ?", tableName, key), fmt.Sprintf("%%%v%%", value))
					}
				}
			} else if paramType == "equal" {
				for key, value := range paramsValue {
					if key != "" {
						db = db.Where(fmt.Sprintf("%v%v = ?", tableName, key), value)
					}
				}
			}
		}
	}

	if p.Distinct {
		err = db.Debug().Select(fmt.Sprintf("count(distinct %sid)", tableName)).Count(&count).Error
	} else {
		err = db.Count(&count).Error
	}

	if !param.NotPage {
		if param.Page < 1 {
			param.Page = 1
		}

		if param.Size < 1 {
			param.Size = 10
		}

		if param.Sort == 0 || param.Sort == -1 {
			orderValue := "id desc"
			if tableName != "" {
				orderValue = tableName + "id desc"
			}
			db = db.Order(orderValue)
		}

		if param.Page == 1 {
			offset = 0
		} else {
			offset = (param.Page - 1) * param.Size
		}
		db = db.Limit(param.Size).Offset(offset)
	}

	err = db.Scan(result).Error
	if err != nil {
		logger.Errorf("数据查询失败，错误：%v", err)
		return nil, err
	}

	paginator.Total = count
	paginator.List = result
	paginator.Page = param.Page
	paginator.Size = param.Size
	paginator.TotalPage = int(math.Ceil(float64(count) / float64(param.Size)))
	if paginator.TotalPage < 0 {
		paginator.TotalPage = 1
	}

	return &paginator, nil
}
