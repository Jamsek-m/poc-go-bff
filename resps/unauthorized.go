package resps

import (
	"encoding/json"
	"net/http"
	"poc-go-bff/oauth2"
	"strconv"
)

func Handle401(res http.ResponseWriter, reason string) {
	res.WriteHeader(http.StatusUnauthorized)
	res.Header().Add(oauth2.HeaderErrReason, reason)
	res.Header().Set("Content-Type", "application/json")

	payload := make(map[string]string)
	payload["status"] = strconv.Itoa(http.StatusUnauthorized)
	payload["message"] = reason

	encoder := json.NewEncoder(res)
	if err := encoder.Encode(payload); err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}
