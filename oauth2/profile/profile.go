package profile

import (
	"fmt"
	"io"
	"net/http"
	"poc-go-bff/config"
	"poc-go-bff/oauth2"
	"poc-go-bff/session"
	"strings"
)

func GetUserProfile(res http.ResponseWriter, req *http.Request) {
	existingSession, _ := session.GetStore().Get(req, config.GetConfig().Sessions.CookieName)
	if existingSession != nil && existingSession.Values["access_token"] != nil {
		request, _ := http.NewRequest("GET", config.GetConfig().Openid.UserinfoURL, nil)

		accessToken := existingSession.Values["access_token"].(string)
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

		client := &http.Client{}
		response, _ := client.Do(request)

		if response.StatusCode == 200 {
			buffer := new(strings.Builder)
			_, _ = io.Copy(buffer, response.Body)

			fmt.Println(buffer.String())
			res.WriteHeader(http.StatusOK)
		} else {
			fmt.Printf("%d: %v", response.StatusCode, response.Status)
			res.Header().Add(oauth2.HeaderErrReason, "Error calling userinfo endpoint!")
			res.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		res.WriteHeader(http.StatusUnauthorized)
		res.Header().Add(oauth2.HeaderErrReason, "No active sessions!")
	}
}
