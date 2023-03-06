package login

import (
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"net/http"
	"net/url"
	"poc-go-bff/config"
	"poc-go-bff/session"
)

func StartCodeFlow(res http.ResponseWriter, req *http.Request) {
	existingSession, err := session.GetStore().Get(req, config.GetConfig().Sessions.CookieName)
	if existingSession != nil && existingSession.Values["access_token"] != nil {
		http.Redirect(res, req, config.GetConfig().Openid.RedirectURL, http.StatusSeeOther)
		return
	}

	newSession, err := session.GetStore().New(req, config.GetConfig().Sessions.CookieName)
	verifier, err := cv.CreateCodeVerifierWithLength(80)
	if err != nil {
		panic("Error creating verifier!")
	}
	newSession.Values["verifier"] = verifier.Value
	newSession.Save(req, res)

	challenge := verifier.CodeChallengeS256()

	params := url.Values{}
	params.Set("response_type", "code")
	params.Set("client_id", "simpl-bff")
	params.Set("redirect_uri", "http://localhost:5000/oidc/callback")
	params.Set("scope", "openid email profile")
	params.Set("code_challenge_method", "S256")
	params.Set("code_challenge", challenge)

	http.Redirect(res, req, "https://auth.gume1a.com/realms/gume1a-int/protocol/openid-connect/auth?"+params.Encode(), 303)
}
