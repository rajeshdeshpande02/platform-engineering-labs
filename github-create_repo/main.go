package main

import (
	"net/http"
	"encoding/json"
	"fmt"
	"bytes"
	"os"

)

const (
	url = "https://api.github.com/user/repos"
)

type RepoDetails struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Private     bool   `json:"private"`
	IsTemplate  bool   `json:"is_template"`
	Auto_init   bool   `json:"auto_init"`
}

func main() {

	github_token := os.Getenv("GHUB_TOKEN")

	if github_token == "" {
		fmt.Println("Error getting GHUB_TOKEN from environment variables")
		return
	}

	RepoDetails := RepoDetails{
		Name:        "test-repo",
		Description: "This is a test repository",
		Homepage:    "https://example.com",
		Private:     false,
		IsTemplate:  false,
		Auto_init:   true,

	}

	repoDetailsJson, err := json.Marshal(RepoDetails)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	// Create a new HTTP client
	httpClient := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(repoDetailsJson))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer " + github_token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")


	//send the request
	resp,  err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()
    if resp.StatusCode != http.StatusCreated {
		fmt.Println("Error creating repository:", resp.Status)
		return	
	}
	fmt.Println("Repository created successfully")


}