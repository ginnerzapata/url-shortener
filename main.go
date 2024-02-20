package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func shortenUrl(url string) string {
	hasher := sha256.New()
	hasher.Write([]byte(url))
	sha := hex.EncodeToString(hasher.Sum(nil))
	return strings.Join([]string{"http://localhost:8080", sha[:8]}, "/")
}

var urlMap = make(map[string]string)

func storeUrl(shortUrl, originalURL string) {
	urlMap[shortUrl] = originalURL
}
func main() {
	http.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			originalURL := r.FormValue("url")
			shortURL := shortenUrl(originalURL)
			storeUrl(shortURL, originalURL)
			fmt.Fprintf(w, "Short URL: %s", shortURL)
		} else {
			http.Error(w, "Invalid reqiest method", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		shortURL := "http://localhost:8080" + r.URL.Path
		fmt.Printf("PATH: %v", shortURL)
		if originalURL, ok := urlMap[shortURL]; ok {
			http.Redirect(w, r, originalURL, http.StatusFound)
		} else {
			http.Error(w, "URL not found", http.StatusNotFound)
		}
	})


	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL Shortener Service")
	})
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}