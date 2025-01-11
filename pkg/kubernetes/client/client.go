package client

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

/*
  @Author : lanyulei
  @Desc :
*/

var (
	config  *rest.Config
	clients Clients
)

type Clients struct {
	clientSet kubernetes.Interface
}

func NewClients() {
	var (
		err error
	)

	// 1. 加载配置，生成配置文件对象。
	config, err = clientcmd.BuildConfigFromFlags("", "/Users/lanyulei/.kube/config")
	if err != nil {
		return
	}

	// 2. 实例化各种客户端
	clients.clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		return
	}
}

func (c *Clients) ClientSet() kubernetes.Interface {
	return c.clientSet
}

func GetConfig() *rest.Config {
	return config
}
