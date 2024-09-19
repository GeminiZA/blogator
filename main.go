package main

import (
	"GeminiZA/blogator/handlers"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")

	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthz", handlers.HandleReadiness)
	mux.HandleFunc("/v1/err", handlers.HandleError)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
