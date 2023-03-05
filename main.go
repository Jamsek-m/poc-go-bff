package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"poc-go-bff/oauth2/login"
	"poc-go-bff/oauth2/profile"
)

func main() {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	// router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/login", login.StartCodeFlow)
	router.Get("/logout", login.LogoutUser)
	router.Get("/profile", profile.GetUserProfile)
	router.Get("/oidc/callback", login.HandleCode)

	err := http.ListenAndServe("0.0.0.0:5000", router)
	if err != nil {
		fmt.Println(err)
	}
}
