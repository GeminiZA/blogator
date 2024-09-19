package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/v1/healthz", handleReadiness)

	server := http.Server{
		Addr:    fmt.Sprintf(":%s", "8080"),
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
