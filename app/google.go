package app

import (
	"context"
	"log"
	"net/http"
	"strings"

	models "whos-that-pokemon/models"
	u "whos-that-pokemon/utils"

	"google.golang.org/api/oauth2/v1"
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
		ctx := r.Context()
		token, err := config.Exchange(ctx, splittedToken[1])
		if err != nil {
			response = u.Message(false, "Something went wrong creating a token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}
		if !token.Valid() {
			// client := config.Client(ctx, token)
		}

		tokenSource := config.TokenSource(ctx, token)
		oauth2Service, err := oauth2.NewService(ctx, option.WithTokenSource(tokenSource))

		if err != nil {
			response = u.Message(false, "Error while authenticating with Google Oauth2.0")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}
		user, err := oauth2Service.Userinfo.Get().Do()
		if err != nil {
			response = u.Message(false, "Error to access User Profile")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			u.Response(w, response)
			return
		}
		account := &models.User{}

		ctx = context.WithValue(ctx, account, user)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)

	})
}
