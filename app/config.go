package app

import (
	"os"

	ouath2 "golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	googleOauth2 "google.golang.org/api/oauth2/v2"
)

var config = &ouath2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{googleOauth2.UserinfoProfileScope, googleOauth2.UserinfoEmailScope},
	Endpoint:     google.Endpoint,
}
