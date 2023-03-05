package profile

import (
	"fmt"
	"io"
	"net/http"
	"poc-go-bff/session"
	"strings"
)

func GetUserProfile(res http.ResponseWriter, req *http.Request) {
	existingSession, _ := session.GetSessionFromCookie("BFF_ID", req)
	if existingSession != nil && existingSession.Tokens != nil {

		fmt.Println("TOKEN: ", existingSession.Tokens.AccessToken)
		request, _ := http.NewRequest("GET", "https://auth.gume1a.com/realms/gume1a-int/protocol/openid-connect/userinfo", nil)
		request.Header.Add("Authorization", "Bearer "+existingSession.Tokens.AccessToken)

		client := &http.Client{}
		response, _ := client.Do(request)

		if response.StatusCode == 200 {
			buffer := new(strings.Builder)
			_, _ = io.Copy(buffer, response.Body)

			fmt.Println(buffer.String())
			res.WriteHeader(http.StatusOK)
		} else {
			fmt.Printf("%d: %v", response.StatusCode, response.Status)
			res.WriteHeader(http.StatusUnauthorized)
		}
	} else {
		res.WriteHeader(http.StatusUnauthorized)
	}
}
