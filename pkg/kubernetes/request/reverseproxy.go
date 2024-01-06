package request

import (
	"bytes"
	"comet/pkg/respstatus"
	"encoding/json"
	"fmt"
	"github.com/lanyulei/toolkit/logger"
	"github.com/lanyulei/toolkit/response"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

/*
  @Author : lanyulei
  @Desc :
*/

func NewSingleHostReverseProxy(target *url.URL) *httputil.ReverseProxy {
	director := func(req *http.Request) {
		rewriteRequestURL(req, target)
	}
	return &httputil.ReverseProxy{Director: director, ModifyResponse: modifyResponseFunc}
}

func rewriteRequestURL(req *http.Request, target *url.URL) {
	targetQuery := target.RawQuery
	req.URL.Scheme = target.Scheme
	req.URL.Host = target.Host
	req.URL.Path, req.URL.RawPath = joinURLPath(target, req.URL)
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
}

func joinURLPath(a, b *url.URL) (path, rawpath string) {
	if a.RawPath == "" && b.RawPath == "" {
		return singleJoiningSlash(a.Path, b.Path), ""
	}
	// Same as singleJoiningSlash, but uses EscapedPath to determine
	// whether a slash should be added
	apath := a.EscapedPath()
	bpath := b.EscapedPath()

	aslash := strings.HasSuffix(apath, "/")
	bslash := strings.HasPrefix(bpath, "/")

	switch {
	case aslash && bslash:
		return a.Path + b.Path[1:], apath + bpath[1:]
	case !aslash && !bslash:
		return a.Path + "/" + b.Path, apath + "/" + bpath
	}
	return a.Path + b.Path, apath + bpath
}

// 重写请求返回格式
func modifyResponseFunc(res *http.Response) error {
	var (
		payloadMap map[string]interface{}
		newPayLoad []byte
	)

	oldPayload, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(oldPayload, &payloadMap)
	if err != nil {
		return err
	}

	payloadResponse := response.Success

	if res.StatusCode >= 200 && res.StatusCode < 300 {
		payloadResponse.Data = payloadMap
	} else {
		// 重写返回格式失败
		logger.Errorf("modifyResponseFunc oldPayload: %s", string(oldPayload))
		payloadResponse.Code = respstatus.UnknownError.Code
		payloadResponse.Message = payloadMap["message"].(string)
	}

	newPayLoad, err = json.Marshal(payloadResponse)
	if err != nil {
		return err
	}

	res.Body = io.NopCloser(bytes.NewBuffer(newPayLoad))
	res.ContentLength = int64(len(newPayLoad))
	res.Header.Set("Content-Length", fmt.Sprint(len(newPayLoad)))
	return nil
}
