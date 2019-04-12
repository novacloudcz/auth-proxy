package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/rs/cors"
)

func main() {
	// testGQL()

	proxyURLString := os.Getenv("PROXY_URL")
	if proxyURLString == "" {
		panic(fmt.Errorf("Missing PROXY_URL environment variable"))
	}
	proxyURL, err := url.Parse(proxyURLString)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	proxy := httputil.NewSingleHostReverseProxy(proxyURL)
	mux.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/", withValidation(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	}))

	handler := cors.AllowAll().Handler(mux)

	log.Fatal(http.ListenAndServe(":"+getEnv("PORT", "80"), handler))
}

// Get env var or default
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
