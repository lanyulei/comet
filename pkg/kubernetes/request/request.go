package request

import (
	"fmt"
	"github.com/lanyulei/comet/pkg/kubernetes/client"
	"github.com/lanyulei/comet/pkg/logger"
	"github.com/lanyulei/comet/pkg/tools/response"
	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/endpoints/filters"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/client-go/rest"
	"net/http"
	"net/url"
)

/*
  @Author : lanyulei
  @Desc :
*/

type responder struct{}

func (r *responder) Error(w http.ResponseWriter, _ *http.Request, err error) {
	logger.Error(err)
	responsewriters.WriteRawJSON(http.StatusOK, response.Response{
		Code:    response.UnknownError.Code,
		Message: fmt.Sprintf("调用 k8s 接口失败，错误：%s", err.Error()),
	}, w)
}

type kubeAPIProxy struct {
	next          http.Handler
	kubeAPIServer *url.URL
	transport     http.RoundTripper
}

// WithKubeAPIServer proxy request to kubernetes service if requests path starts with /api
func WithKubeAPIServer(next http.Handler, config *rest.Config) http.Handler {
	kubeAPIServer, _ := url.Parse(config.Host)
	transport, err := rest.TransportFor(config)
	if err != nil {
		logger.Errorf("Unable to create transport from rest.Config: %v", err)
		return next
	}
	return &kubeAPIProxy{
		next:          next,
		kubeAPIServer: kubeAPIServer,
		transport:     transport,
	}
}

func (k kubeAPIProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s := *req.URL
	s.Host = k.kubeAPIServer.Host
	s.Scheme = k.kubeAPIServer.Scheme

	// make sure we don't override kubernetes's authorization
	req.Header.Del("Authorization")
	httpProxy := proxy.NewUpgradeAwareHandler(&s, k.transport, true, false, &responder{})
	httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(k.transport, k.transport)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, HEAD")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Length")
	w.Header().Set("Access-Control-Max-Age", "43200")

	if req.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	httpProxy.ServeHTTP(w, req)
}

func BuildHandlerChain(h http.Handler) (handler http.Handler) {
	// k8s 接口转发
	requestInfoResolver := &request.RequestInfoFactory{
		APIPrefixes:          sets.NewString("apis", "api"),
		GrouplessAPIPrefixes: sets.NewString("apis", "api"),
	}

	handler = WithKubeAPIServer(h, client.GetConfig())
	handler = filters.WithRequestInfo(handler, requestInfoResolver)
	return
}
