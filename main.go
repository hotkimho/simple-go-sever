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

	if corsOrigin == "" {
		corsOrigin = "*(default value)"
	}
	if debugMode == "" {
		debugMode = "false(default value)"
	}
	if logLevel == "" {
		logLevel = "info(default value)"
	}

	fmt.Println("trigger test-2")
	fmt.Printf("PORT: %s\n", port)
	fmt.Printf("CORS_ORIGIN: %s\n", corsOrigin)
	fmt.Printf("DEBUG_MODE: %s\n", debugMode)
	fmt.Printf("LOG_LEVEL: %s\n", logLevel)

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Test 3")
		// content-type값이 application/json으로 설정하고 chartset=utf-8 설정 삭제
		w.Header().Del("Content-Type")
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
	})

	log.Printf("Starting server on port %s...", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
