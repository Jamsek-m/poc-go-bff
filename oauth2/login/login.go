package login

import (
	"fmt"
	cv "github.com/nirasan/go-oauth-pkce-code-verifier"
	"net/http"
	"net/url"
	"poc-go-bff/session"
	"poc-go-bff/session/store"
)

func StartCodeFlow(res http.ResponseWriter, req *http.Request) {
	existingSession, err := session.GetSessionFromCookie("BFF_ID", req)
	if existingSession != nil && existingSession.Tokens != nil {
		http.Redirect(res, req, "https://google.com", http.StatusSeeOther)
		return
	}

	sessionID, newSession := session.GetStore().Start()
	sessionCookie := session.NewSessionCookie(sessionID, res)
	fmt.Println(sessionCookie)
	http.SetCookie(res, sessionCookie)

	verifier, err := cv.CreateCodeVerifierWithLength(80)
	if err != nil {
		panic("Error creating verifier!")
	}
	newSession.CodeFlow = &store.CodeFlow{
		Verifier: verifier.Value,
	}
	session.GetStore().Update(sessionID, newSession)

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
