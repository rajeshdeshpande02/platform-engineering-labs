package main

import (
	"github-create_branch/internal/service"
	"net/http"
	"encoding/json"
	"log"
)



func main(){
	mux := http.NewServeMux()

	mux.HandleFunc("/create-branch", createBranchHandler) 

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Starting server on :8080")

	if err := server.ListenAndServe(); err != nil  {
		log.Fatalf("Error starting server: %v", err)
	}
	log.Println("Server stopped")    

}

func createBranchHandler(w http.ResponseWriter, r *http.Request) {
	// Declaring it locally to  make it thread sae
	var requestData struct {
		RepoName   string `json:"repo-name"`
		BaseBranch string `json:"base-branch"`
		NewBranch  string `json:"new-branch"`
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

	repoName := requestData.RepoName
	baseBranch := requestData.BaseBranch
	newBranch := requestData.NewBranch

	if repoName == "" || baseBranch == "" || newBranch == "" {
		http.Error(w, "repo-name, base-branch and new-branch query parameters are required", http.StatusBadRequest)
		return
	}

	err = service.CreateBranch(repoName, baseBranch, newBranch)
	if err != nil {
		http.Error(w, "Error creating branch: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Branch created successfully"))

}