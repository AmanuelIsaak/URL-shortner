package main

import (
	"errors"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
)

type PageData struct {
	Short string
}

type shortURL struct {
	OriginalURL string `json:"original_url"`
	ShortCode   string `json:"short_code"`
}

var urlMap = make(map[string]shortURL)

func main() {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "index.html", PageData{})
		if err != nil {
			return
		}
	})

	router.HandleFunc("/shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			originalURL := r.FormValue("url")
			shortCode := generateShortCode()
			urlMap[shortCode] = shortURL{OriginalURL: originalURL, ShortCode: shortCode}

			err := tmpl.ExecuteTemplate(w, "shorten.html", PageData{
				Short: "http:/localhost:8080/r/" + shortCode,
			})
			if err != nil {
				return
			}
		} else {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		}
	})

	router.HandleFunc("/r/", func(w http.ResponseWriter, r *http.Request) {
		shortCode := r.URL.Path[len("/r/"):]
		if url, ok := urlMap[shortCode]; ok {
			http.Redirect(w, r, url.OriginalURL, http.StatusFound)
		} else {
			http.NotFound(w, r)
		}
	})

	srv := http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	fmt.Println("Starting website at localhost:8080")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		fmt.Println("An error occured:", err)
	}
}

func generateShortCode() string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortCode := make([]byte, 7)
	for i := range shortCode {
		shortCode[i] = chars[rand.Intn(len(chars))]
	}
	return string(shortCode)
}
