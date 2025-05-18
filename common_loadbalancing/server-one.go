package main

import (
	"log"
	"net/http"

)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serverHandler) 

	server := &http.Server{
		Addr:   ":8080",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from Server One!"))
	log.Println("Server One: Request received")
	log.Println("Server One: Response sent")

}
