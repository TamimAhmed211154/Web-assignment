package main

import (
	"fmt"
	"net/http"
)

var urls = map[string]string{}

func shorten(w http.ResponseWriter, r *http.Request) {
	longURL := r.URL.Query().Get("url")

	if longURL == "" {
		fmt.Fprintln(w, "URL is required")
		return
	}

	code := "abc123"
	urls[code] = longURL

	fmt.Fprintln(w, "Short URL: http://localhost:8080/"+code)
}

func redirect(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path[1:]

	if code == "" {
		fmt.Fprintln(w, "URL Shortener is running")
		return
	}

	longURL := urls[code]

	if longURL == "" {
		fmt.Fprintln(w, "Short URL not found")
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", shorten)
	http.HandleFunc("/", redirect)

	fmt.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
}


//http://localhost:8080/shorten?url=https://www.google.com