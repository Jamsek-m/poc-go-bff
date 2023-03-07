package login

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"poc-go-bff/config"
	"poc-go-bff/oauth2"
	"poc-go-bff/session"
	"strings"
)

func HandleCode(res http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get(oauth2.CodeQueryParam)
	if code == "" {
		res.WriteHeader(http.StatusUnauthorized)
		res.Header().Add(oauth2.HeaderErrReason, "No code query parameter present!")
		return
	}

	verifier := session.Current().GetSessionValue(req, "verifier")
	if verifier == nil {
		res.WriteHeader(http.StatusUnauthorized)
		res.Header().Add(oauth2.HeaderErrReason, "No active session that started login flow!")
		return
	}

	params := url.Values{}
	params.Set("grant_type", "authorization_code")
	params.Set("code", code)
	params.Set("redirect_uri", config.GetConfig().Openid.CallbackURL)
	params.Set("scopes", config.GetConfig().Openid.Scopes)
	params.Set("code_verifier", verifier.(string))
	params.Set("client_id", config.GetConfig().Openid.ClientID)
	params.Set("client_secret", config.GetConfig().Openid.ClientSecret)

	request, _ := http.NewRequest("POST", config.GetConfig().Openid.TokenURL, strings.NewReader(params.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, _ := client.Do(request)

	if response.StatusCode == 200 {
		buffer := new(strings.Builder)
		_, _ = io.Copy(buffer, response.Body)

		data := oauth2.TokenResponse{}
		_ = json.Unmarshal([]byte(buffer.String()), &data)

		session.Current().SetSessionValue(req, "access_token", data.AccessToken)
		session.Current().SetSessionValue(req, "refresh_token", data.RefreshToken)
		session.Current().SetSessionValue(req, "id_token", data.IDToken)
		session.Current().CommitSessionChanges(req)

		http.Redirect(res, req, config.GetConfig().Openid.RedirectURL, http.StatusSeeOther)
	} else {
		fmt.Printf("%d: %v", response.StatusCode, response.Status)
		res.Header().Add(oauth2.HeaderErrReason, "Error when retrieving tokens!")
		res.WriteHeader(http.StatusUnauthorized)
	}
}
