package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {

	proxyURL := getEnvURL("PROXY_URL")
	jwksProviderURL := getEnv("JWKS_PROVIDER_URL")
	requiredJWTScopes := getEnvWithFallback("REQUIRED_JWT_SCOPES", "")

	mux := http.NewServeMux()

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	scopesArray := []string{}
	if requiredJWTScopes != "" {
		scopesArray = strings.Split(requiredJWTScopes, " ")
	}
	vOptions := ValidationOptions{
		jwksProviderURL:   jwksProviderURL,
		requiredJWTScopes: scopesArray,
	}
	mux.HandleFunc("/", withValidation(func(w http.ResponseWriter, r *http.Request) {
		r.Header.Del("authorization")
		q := r.URL.Query()
		q.Del("access_token")
		r.URL.RawQuery = q.Encode()
		r.Host = proxyURL.Host
		proxy.ServeHTTP(w, r)
	}, vOptions))

	port := getEnvWithFallback("PORT", "80")
	log.Printf("running at http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}

// Get env var or default
func getEnvWithFallback(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(fmt.Errorf("Missing %s environment variable", key))
}

func getEnvURL(key string) *url.URL {
	url, err := url.Parse(getEnv(key))
	if err != nil {
		panic(err)
	}
	return url
}
