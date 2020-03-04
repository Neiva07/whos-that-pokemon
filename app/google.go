package app

import (
	"context"
	"log"
	"net/http"
	"strings"
	models "whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

	"golang.org/x/oauth2"
	googleOauth "google.golang.org/api/oauth2/v1"
	"google.golang.org/api/option"
)

//Authentication layer to valid requests from the client
var Authentication = func(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		notAuth := []string{"/api/users/signin"}

		requestPath := r.URL.Path
		log.Println(requestPath)

		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		response := make(map[string]interface{})
		tokenHeader := r.Header.Get("Authorization")

		if tokenHeader == "" {
			response = u.Message(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}

		splittedToken := strings.Split(tokenHeader, " ")

		if len(splittedToken) != 2 {
			response = u.Message(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}
		authCode := splittedToken[1]

		user := &models.User{}

		err := models.Redis.Get(authCode).Scan(user)

		if err != nil {
			response = u.Message(false, "Error retrieving user from the session system.")
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Add("Content-Type", "application/json")
			response["error"] = err
			log.Println(err)
			u.Response(w, response)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}

//TokenHandler receive authCode to exchange for the respective token
var TokenHandler = func(authCode string, ctx context.Context) (*oauth2.Token, error) {

	token, err := config.Exchange(ctx, authCode)

	return token, err
}

//GetGoogleUserInfo returns the user from the google profile api
var GetGoogleUserInfo = func(token *oauth2.Token, ctx context.Context) (*googleOauth.Userinfoplus, error) {
	tokenSource := config.TokenSource(ctx, token)
	oauth2Service, err := googleOauth.NewService(ctx, option.WithTokenSource(tokenSource))

	if err != nil {
		return nil, err
	}
	user, err := oauth2Service.Userinfo.Get().Do()

	return user, err
}
