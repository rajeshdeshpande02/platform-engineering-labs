package main

import (
	"net/http"
	"encoding/json"
	"log"
	"github-create_repo_from_template/internal/service"

)

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/repo-with-template", createRepoHandler)

	server := &http.Server {
		Addr:    ":8080",
		Handler: mux,
	}

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Println("Server stopped")
}

func createRepoHandler(w http.ResponseWriter, r *http.Request) {

	var requestData struct {
		Template string `json:"template"`
		Repo_Name string `json:"repo_name"`
	}

	if r.Method != http.MethodPost {
		log.Printf("Error Invalid request method: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}
	template := requestData.Template
	repo_name := requestData.Repo_Name
	if template == "" || repo_name == "" {	
		http.Error(w, "template and repo_name are required", http.StatusBadRequest)
		return
	}
	err = service.CreateRepoFromTempl(template, repo_name)
	if err != nil {
		http.Error(w, "Error creating repository from template: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Write([]byte("Repository created successfully"))

	w.WriteHeader(http.StatusOK)
}