package login

import (
	"fmt"
	"net/http"
	"net/url"
	"poc-go-bff/config"
	"poc-go-bff/session"
)

func LogoutUser(res http.ResponseWriter, req *http.Request) {
	existingSession, _ := session.GetStore().Get(req, config.GetConfig().Sessions.CookieName)
	if existingSession != nil && existingSession.Values["access_token"] != nil {

		idToken := existingSession.Values["id_token"].(string)

		params := url.Values{}
		params.Set("id_token_hint", idToken)
		params.Set("post_logout_redirect_uri", config.GetConfig().Openid.RedirectURL)

		// Destroys session
		existingSession.Options.MaxAge = -1
		existingSession.Save(req, res)

		endSessionUrl := fmt.Sprintf("%s?%s", config.GetConfig().Openid.EndSessionURL, params.Encode())
		http.Redirect(res, req, endSessionUrl, 303)
	} else {
		http.Redirect(res, req, config.GetConfig().Openid.RedirectURL, 303)
	}
}
