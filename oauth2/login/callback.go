package login

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"poc-go-bff/oauth2"
	"poc-go-bff/session"
	"poc-go-bff/session/store"
	"strings"
)

func HandleCode(res http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	if code == "" {
		fmt.Println("Code not found!")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	existingSession, _ := session.GetSessionFromCookie("BFF_ID", req)
	if existingSession == nil || existingSession.CodeFlow == nil {
		fmt.Println("Session/cookie not found!")
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("Preq fulfilled, continue...")
	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("redirect_uri", "http://localhost:5000/oidc/callback")
	params.Set("scopes", "openid email profile")
	params.Set("code_verifier", existingSession.CodeFlow.Verifier)
	params.Set("client_id", "simpl-bff")
	params.Set("client_secret", "H6srj8ArUZtL0r0ju10Pg3FEQ4MOd9DA")

	request, _ := http.NewRequest("POST", "https://auth.gume1a.com/realms/gume1a-int/protocol/openid-connect/token", strings.NewReader(params.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, _ := client.Do(request)

	if response.StatusCode == 200 {
		buffer := new(strings.Builder)
		_, _ = io.Copy(buffer, response.Body)

		data := oauth2.TokenResponse{}
		_ = json.Unmarshal([]byte(buffer.String()), &data)

		fmt.Println(data)
		existingSession.Tokens = &store.UserTokens{
			AccessToken:  data.AccessToken,
			RefreshToken: data.RefreshToken,
			IDToken:      data.IDToken,
		}
		session.GetStore().Update(*existingSession.ID, existingSession)

		http.Redirect(res, req, "https://google.com", http.StatusSeeOther)
	} else {
		fmt.Printf("%d: %v", response.StatusCode, response.Status)
		res.WriteHeader(http.StatusUnauthorized)
	}
}
