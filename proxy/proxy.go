package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"poc-go-bff/oauth2"
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
		res.WriteHeader(http.StatusUnauthorized)
		res.Header().Add(oauth2.HeaderErrReason, "No active sessions!")
	}
	p.proxy.ServeHTTP(res, req)
}

func New(URL *url.URL) ApiProxy {
	fmt.Println("New: ", URL)
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
