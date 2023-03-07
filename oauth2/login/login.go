package login

import (
	"fmt"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"net/http"
	"net/url"
	"poc-go-bff/config"
	"poc-go-bff/session"
)

func StartCodeFlow(res http.ResponseWriter, req *http.Request) {
	accessToken := session.Current().GetSessionValue(req, "access_token")
	if accessToken != nil {
		http.Redirect(res, req, config.GetConfig().Openid.RedirectURL, http.StatusSeeOther)
		return
	}
	session.Current().CreateSession(req)
	verifier, err := cv.CreateCodeVerifierWithLength(80)
	if err != nil {
		panic("Error creating verifier!")
	}

	session.Current().SetSessionValue(req, "verifier", verifier.Value)
	session.Current().CommitSessionChanges(req)

	challenge := verifier.CodeChallengeS256()

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", config.GetConfig().Openid.ClientID)
	params.Set("redirect_uri", config.GetConfig().Openid.CallbackURL)
	params.Set("scope", config.GetConfig().Openid.Scopes)
	params.Set("code_challenge_method", *config.GetConfig().Openid.PKCEMethod)
	params.Set("code_challenge", challenge)

	authorizationURL := fmt.Sprintf("%s?%s", config.GetConfig().Openid.AuthorizationURL, params.Encode())
	http.Redirect(res, req, authorizationURL, http.StatusSeeOther)
}
