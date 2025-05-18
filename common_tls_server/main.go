package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/greet", greetHandler)
	server := &http.Server{
		Addr:    ":6443",
		Handler: mux,
	}
	log.Println("Starting server on :6443")

	err := server.ListenAndServeTLS("server.crt", "server.key")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	log.Println("Server stopped")
}

func greetHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for /greet")
	w.Write([]byte("Hello, World!"))
}
