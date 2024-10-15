package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	corsOrigin := os.Getenv("CORS_ORIGIN")
	debugMode := os.Getenv("DEBUG_MODE")
	logLevel := os.Getenv("LOG_LEVEL")

	fmt.Printf("PORT: %s\n", port)
	fmt.Printf("CORS_ORIGIN: %s\n", corsOrigin)
	fmt.Printf("DEBUG_MODE: %s\n", debugMode)
	fmt.Printf("LOG_LEVEL: %s\n", logLevel)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test")
	})

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
