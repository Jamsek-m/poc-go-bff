package login

import (
	"net/http"
	"net/url"
	"poc-go-bff/session"
)

func LogoutUser(res http.ResponseWriter, req *http.Request) {
	existingSession, _ := session.GetSessionFromCookie("BFF_ID", req)
	if existingSession != nil && existingSession.Tokens != nil {
		params := url.Values{}
		params.Set("id_token_hint", existingSession.Tokens.IDToken)
		params.Set("post_logout_redirect_uri", "https://google.com")

		session.GetStore().Destroy(*existingSession.ID)

		http.Redirect(res, req, "https://auth.gume1a.com/realms/gume1a-int/protocol/openid-connect/logout?"+params.Encode(), 303)
	} else {
		http.Redirect(res, req, "https://google.com", 303)
	}
}
