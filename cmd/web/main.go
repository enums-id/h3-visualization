package main

import (
	"log"
	"net/http"

	"h3-visualization/internal/web"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", web.IndexHandler)
	mux.HandleFunc("POST /run", web.RunHandler)
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Println("H3 Visualization server listening on http://localhost:8080")
	log.Println("Run this command from the repository root directory.")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
