package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"

	"github.com/jakubknejzlik/go-jwks"
)

func withValidation(next http.HandlerFunc) http.HandlerFunc {
	client, err := jwks.NewClient("https://id.novacloud.cz/.well-known/jwks.json")
	if err != nil {
		panic(err)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromRequest(r)

		if token == "" || validateToken(client, token) != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Header().Set("content-type", "text/plain")
			fmt.Fprintf(w, "401 Unauthorized")
			return
		}

		next.ServeHTTP(w, r)
	}
}

func tokenFromRequest(r *http.Request) (token string) {
	token = r.Header.Get("authorization")
	if token == "" {
		token = r.URL.Query().Get("access_token")
	}

	token = strings.Replace(token, "Bearer ", "", 1)

	return
}

func validateToken(client *jwks.Client, token string) (err error) {
	keys, err := client.GetKeys()
	if err != nil {
		return
	}
	for _, key := range keys {
		_, err = jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}
			return key, nil
		})
		if err == nil {
			return
		}
	}

	return
}
