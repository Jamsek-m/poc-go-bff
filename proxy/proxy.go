package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"poc-go-bff/resps"
	"poc-go-bff/session"
)

type (
	apiProxy struct {
		proxy *httputil.ReverseProxy
	}

	ApiProxy interface {
		ProxyHttp(res http.ResponseWriter, req *http.Request)
	}
)

func (p *apiProxy) ProxyHttp(res http.ResponseWriter, req *http.Request) {
	req.URL.Path = ""
	req.RequestURI = ""

	accessToken := session.Current().GetSessionValue(req, "access_token")
	if accessToken != nil {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken.(string)))
	} else {
		resps.Handle401(res, "No active sessions!")
	}
	p.proxy.ServeHTTP(res, req)
}

func New(URL *url.URL) ApiProxy {
	proxy := httputil.NewSingleHostReverseProxy(URL)
	defaultDirector := proxy.Director

	proxy.Director = func(req *http.Request) {
		defaultDirector(req)
		req.Header.Add("X-Proxy", "BFF_PROXY")
		req.Host = URL.Host
	}

	return &apiProxy{
		proxy: proxy,
	}
}
