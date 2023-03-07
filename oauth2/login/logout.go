package login

import (
	"fmt"
	"net/http"
	"net/url"
	"poc-go-bff/config"
	"poc-go-bff/session"
)

func LogoutUser(res http.ResponseWriter, req *http.Request) {

	idToken := session.Current().GetSessionValue(req, "id_token")
	if idToken != nil {
		params := url.Values{}
		params.Set("id_token_hint", idToken.(string))
		params.Set("post_logout_redirect_uri", config.GetConfig().Openid.RedirectURL)

		if err := session.Current().DestroySession(req); err == nil {
			fmt.Println(err)
		}
		session.Current().CommitSessionChanges(req)

		endSessionUrl := fmt.Sprintf("%s?%s", config.GetConfig().Openid.EndSessionURL, params.Encode())
		http.Redirect(res, req, endSessionUrl, http.StatusSeeOther)
	} else {
		http.Redirect(res, req, config.GetConfig().Openid.RedirectURL, http.StatusSeeOther)
	}
}
