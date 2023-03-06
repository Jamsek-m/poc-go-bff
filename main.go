package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/gorilla/mux"
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

	session.InitSessions()

	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/login", login.StartCodeFlow).Methods("GET")
	router.HandleFunc("/login", login.StartCodeFlow).Methods("GET")
	router.HandleFunc("/logout", login.LogoutUser).Methods("GET")
	router.HandleFunc("/profile", profile.GetUserProfile).Methods("GET")
	router.HandleFunc(config.GetConfig().Openid.CallbackPath, login.HandleCode).Methods("GET")

	addr := fmt.Sprintf("%s:%d", config.GetConfig().Host, config.GetConfig().Port)
	log.Fatal(http.ListenAndServe(addr, router))
}
