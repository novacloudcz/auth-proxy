package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/cors"

	"github.com/jakubknejzlik/go-jwks"
)

type JWTTokenClaims struct {
	jwt.StandardClaims
	Scope string `json:"scope,omitempty"`
}

// ValidationOptions ...
type ValidationOptions struct {
	jwksProviderURL   string //"https://id.novacloud.cz/.well-known/jwks.json"
	requiredJWTScopes []string
}

func withValidation(next http.HandlerFunc, options ValidationOptions) http.HandlerFunc {
	client, err := jwks.NewClient(options.jwksProviderURL)
	if err != nil {
		panic(err)
	}
	corsHandler := cors.AllowAll()
	validateScopes := len(options.requiredJWTScopes) > 0

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != "OPTIONS" {

			token := tokenFromRequest(r)

			valid, claims, err := validateToken(client, token)
			if token == "" || !valid || err != nil {
				corsHandler.ServeHTTP(w, r, func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusUnauthorized)
					w.Header().Set("content-type", "text/plain")
					fmt.Fprintf(w, "401 Unauthorized")
				})
				return
			}

			if os.Getenv("DEBUG") != "" {
				log.Println("Parsed token with subject", claims.Subject, claims.Issuer)
			}

			if validateScopes {
				if !validateRequiredJWTScopes(claims, options.requiredJWTScopes) {
					w.WriteHeader(http.StatusForbidden)
					w.Header().Set("content-type", "text/plain")
					fmt.Fprintf(w, "403 Missing required scope(s)")
					return
				}
			}
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

func validateToken(client *jwks.Client, token string) (valid bool, claims JWTTokenClaims, err error) {
	if token == "" {
		return
	}
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok || kid == "" {
			return nil, fmt.Errorf("Token is missing kid field in headers")
		}
		key, err := client.GetKey(kid)
		if err != nil {
			return nil, err
		}
		return key, nil
	})
	if err != nil {
		return
	}
	if jwtToken != nil {
		valid = jwtToken.Valid
	}

	return
}

func validateRequiredJWTScopes(claims JWTTokenClaims, scopes []string) bool {
	tokenScopes := map[string]struct{}{}

	for _, s := range strings.Split(claims.Scope, " ") {
		tokenScopes[s] = struct{}{}
	}

	for _, s := range scopes {
		_, contains := tokenScopes[s]
		if !contains {
			return false
		}
	}

	return true
}
