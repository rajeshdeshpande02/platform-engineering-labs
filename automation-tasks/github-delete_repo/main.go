package main

import (
	"fmt"
	"net/http"
	"os"
)

const (
	url = "https://api.github.com/user/repos"
)

func main() {

	github_token := os.Getenv("GHUB_TOKEN")

	if github_token == "" {
		fmt.Println("Error getting GHUB_TOKEN from environment variables")
		return
	}
	
	repoName := "rajeshdeshpande02/test-repo"
	url := fmt.Sprintf("https://api.github.com/repos/%s", repoName)
	fmt.Println("URL:", url)

	// Create a new HTTP client
	httpClient := &http.Client{}

	// Create a new request
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Set the request headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer " + github_token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	//send the request
	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println("Error sending delete request:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		fmt.Println("Error creating repository:", resp.Status)
		return
	}
	fmt.Println("Repository delete successfully: ", repoName)

}
