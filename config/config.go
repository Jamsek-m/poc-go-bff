package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	DefaultPort = "5000"
	DefaultHost = "localhost"
)

type (
	AppConfig struct {
		Host    string
		Port    int
		BaseURL string
		Openid  struct {
			ClientID         string
			ClientSecret     string
			AuthorizationURL string
			TokenURL         string
			UserinfoURL      string
			EndSessionURL    string
			RedirectURL      string
			Scopes           string
			PKCEMethod       *string
			CallbackPath     string
			CallbackURL      string
		}
		Sessions struct {
			Store      string
			Secret     string
			CookieName string
		}
	}
)

var config *AppConfig

func getOrDefault(envVar string, defaultValue string) string {
	if value, ok := os.LookupEnv(envVar); ok {
		return value
	}
	return defaultValue
}

func InitConfig() {
	config = &AppConfig{}

	config.Host = getOrDefault("BFF_HOST", DefaultHost)
	stringPort := getOrDefault("BFF_PORT", DefaultPort)
	port, err := strconv.Atoi(stringPort)
	if err != nil {
		panic(err)
	} else {
		config.Port = port
	}

	baseURL, baseURLIsSet := os.LookupEnv("BFF_OIDC_BASE_URL")
	if baseURLIsSet {
		config.BaseURL = baseURL
	} else {
		config.BaseURL = fmt.Sprintf("http://%s:%d", config.Host, config.Port)
	}

	config.Openid.CallbackPath = "/oidc/callback"
	config.Openid.CallbackURL = fmt.Sprintf("%s%s", config.BaseURL, config.Openid.CallbackPath)

	config.Openid.ClientID = os.Getenv("BFF_OIDC_CLIENT_ID")
	config.Openid.ClientSecret = os.Getenv("BFF_OIDC_CLIENT_SECRET")
	config.Openid.AuthorizationURL = os.Getenv("BFF_OIDC_AUTH_URL")
	config.Openid.TokenURL = os.Getenv("BFF_OIDC_TOKEN_URL")
	config.Openid.UserinfoURL = os.Getenv("BFF_OIDC_USERINFO_URL")
	config.Openid.EndSessionURL = os.Getenv("BFF_OIDC_END_SESSION_URL")
	config.Openid.RedirectURL = os.Getenv("BFF_OIDC_REDIRECT_URL")
	config.Openid.Scopes = os.Getenv("BFF_OIDC_SCOPES")
	if pkce, ok := os.LookupEnv("BFF_OIDC_PKCE"); ok {
		if pkce == "S256" || pkce == "plain" {
			config.Openid.PKCEMethod = &pkce
		} else {
			panic("invalid PKCE value")
		}
	}
	store, storeIsSet := os.LookupEnv("BFF_SESSIONS_STORE")
	if !storeIsSet || (store != "memory" && store != "sqlite3" && store != "redis") {
		config.Sessions.Store = "sqlite3"
	} else {
		config.Sessions.Store = store
	}
	config.Sessions.Secret = os.Getenv("BFF_SESSIONS_SECRET")
	config.Sessions.CookieName = "BFF_ID"
}

func GetConfig() *AppConfig {
	return config
}
