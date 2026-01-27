package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// Test simple route
	mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple route works!")
	})

	// Test route with path parameter
	mux.HandleFunc("/test/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		fmt.Fprintf(w, "ID route works! ID: %s", id)
	})

	// Test API route
	mux.HandleFunc("/api/brands", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "API brands route works!")
	})

	log.Println("Test server starting on :8080")
	log.Println("Routes:")
	log.Println("  - /test")
	log.Println("  - /test/{id}")
	log.Println("  - /api/brands")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
