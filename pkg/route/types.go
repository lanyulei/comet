package route

/*
  @Author : lanyulei
  @Desc :
*/

const (
	CheckRouteRegisterPath  = "/api/v1/route/register/check"
	BatchRegisterRouterPath = "/api/v1/system/api/batch"
	Unregistered            = "unregistered"
	Invalid                 = "invalid"
)

type Route struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}
