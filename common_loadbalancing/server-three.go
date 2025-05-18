package main

import (
	"log"
	"net/http"

)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serverHandler) 

	server := &http.Server{
		Addr:   ":8082",
		Handler: mux,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func serverHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello from Server Three!"))
	log.Println("Server Three: Request received")
	log.Println("Server Three: Response sent")

}
