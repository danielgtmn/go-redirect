package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func main() {
	targetURL := os.Getenv("REDIRECT_TARGET")
	if targetURL == "" {
		log.Fatal("REDIRECT_TARGET environment variable is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	redirectCode := 301
	if code := os.Getenv("REDIRECT_CODE"); code != "" {
		if c, err := strconv.Atoi(code); err == nil && (c == 301 || c == 302) {
			redirectCode = c
		}
	}

	preservePath := os.Getenv("PRESERVE_PATH") == "true"

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		target := targetURL
		if preservePath {
			target = targetURL + r.URL.Path
			if r.URL.RawQuery != "" {
				target += "?" + r.URL.RawQuery
			}
		}

		log.Printf("%s %s -> %s (%d)", r.Method, r.URL.String(), target, redirectCode)
		http.Redirect(w, r, target, redirectCode)
	})

	log.Printf("Starting redirect server on :%s", port)
	log.Printf("Target: %s (Code: %d, Preserve Path: %v)", targetURL, redirectCode, preservePath)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
