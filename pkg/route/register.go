package route

import (
	"github.com/lanyulei/comet/pkg/logger"

	"github.com/spf13/viper"

	"github.com/guonaihong/gout"
)

/*
  @Author : lanyulei
  @Desc :
*/

func CheckRegisterRoute(routes []*Route) (err error) {
	var (
		result  map[string]interface{}
		apiInfo []map[string]interface{}
	)

	err = gout.POST(viper.GetString("nebula.url") + CheckRouteRegisterPath).
		SetJSON(gout.H{
			"routes": routes,
		}).
		BindJSON(&result).
		Do()
	if err != nil {
		return
	}

	unregisteredRoutes := result["data"].(map[string]interface{})[Unregistered]
	if unregisteredRoutes != nil && len(unregisteredRoutes.([]interface{})) > 0 {
		for _, r := range unregisteredRoutes.([]interface{}) {
			apiInfo = append(apiInfo, map[string]interface{}{
				"method": r.(map[string]interface{})["method"],
				"url":    r.(map[string]interface{})["path"],
			})
		}
	}

	if len(apiInfo) > 0 {
		err = gout.POST(viper.GetString("nebula.url") + BatchRegisterRouterPath).
			SetJSON(apiInfo).
			BindJSON(&result).
			Do()
		if err != nil {
			for _, v := range apiInfo {
				logger.Warnf("unregistered route: %s %s", v["method"], v["url"])
			}
			return
		}
	}

	return
}
