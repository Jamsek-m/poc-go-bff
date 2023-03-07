package proxy

import (
	"net/url"
	"poc-go-bff/config"
	"sync"
)

var proxy ApiProxy
var lock sync.Mutex

func UserInfoProxy() ApiProxy {
	if proxy == nil {
		lock.Lock()
		defer lock.Unlock()
		if proxy != nil {
			return proxy
		}
		userInfoURL, err := url.Parse(config.GetConfig().Openid.UserinfoURL)
		if err != nil {
			panic(err)
		}
		proxy = New(userInfoURL)
	}
	return proxy
}
