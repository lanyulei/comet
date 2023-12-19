package response

var (
	Success = Response{Code: 20000, Message: "success"}

	AuthorizationNullError   = Response{Code: 30001, Message: "请求头中 Authorization 为空"}
	AuthorizationFormatError = Response{Code: 30002, Message: "请求头中 Authorization 格式有误"}
	InvalidTokenError        = Response{Code: 30003, Message: "Token 无效"}
	UnknownError             = Response{Code: 30004, Message: "未知错误"}
	InvalidParameterError    = Response{Code: 30005, Message: "参数有误"}

	GetResourceListError           = Response{Code: 40001, Message: "获取资源列表失败"}
	GetResourceDetailsError        = Response{Code: 40002, Message: "获取资源详情失败"}
	ClusterConfigListError         = Response{Code: 40003, Message: "获取集群配置列表失败"}
	GetClusterConfigError          = Response{Code: 40004, Message: "获取集群配置详情失败"}
	ClusterConfigExistError        = Response{Code: 40005, Message: "集群配置已存在"}
	ClusterConfigCreateError       = Response{Code: 40006, Message: "创建集群配置失败"}
	ClusterConfigEditError         = Response{Code: 40007, Message: "编辑集群配置失败"}
	ClusterConfigDeleteError       = Response{Code: 40008, Message: "删除集群配置失败"}
	ClusterConfigDefaultExistError = Response{Code: 40009, Message: "集群配置已存在默认配置"}
	SwitchClusterError             = Response{Code: 40010, Message: "切换集群失败"}
	WriteKubeconfigError           = Response{Code: 40011, Message: "写入 kubeconfig 失败"}
	UpdateClusterConfigError       = Response{Code: 40012, Message: "更新集群配置失败"}
	ReadKubeconfigError            = Response{Code: 40013, Message: "读取 kubeconfig 失败"}
	GetUserClusterError            = Response{Code: 40014, Message: "获取用户集群失败"}
	CreateUserClusterError         = Response{Code: 40015, Message: "创建用户集群失败"}
	UpdateUserClusterError         = Response{Code: 40016, Message: "更新用户集群失败"}
	ClusterNotExistError           = Response{Code: 40017, Message: "集群不存在"}
)
