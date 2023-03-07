package profile

import (
	"net/http"
	"poc-go-bff/proxy"
)

func GetUserProfile(res http.ResponseWriter, req *http.Request) {
	proxy.UserInfoProxy().ProxyHttp(res, req)
}
