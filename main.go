package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"poc-go-bff/config"
	"poc-go-bff/oauth2/login"
	"poc-go-bff/oauth2/profile"
	"poc-go-bff/session"
)

func main() {
	err := godotenv.Load()
	if err == nil {
		log.Printf(color.GreenString(".env loaded"))
	} else {
		log.Printf(color.RedString("error loading .env!"))
	}
	config.InitConfig()

	session.InitStore()

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)

	router.Get("/login", login.StartCodeFlow)
	router.Get("/login", login.StartCodeFlow)
	router.Get("/logout", login.LogoutUser)
	router.Get("/profile", profile.GetUserProfile)
	router.Get(config.GetConfig().Openid.CallbackPath, login.HandleCode)

	addr := fmt.Sprintf("%s:%d", config.GetConfig().Host, config.GetConfig().Port)
	sessionHandler := session.Current().Instance().LoadAndSave(router)
	log.Fatal(http.ListenAndServe(addr, sessionHandler))
}
