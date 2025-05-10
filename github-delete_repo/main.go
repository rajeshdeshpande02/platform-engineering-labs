package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type GitHubClient struct {
	HTTPClient *http.Client
	Token      string
}

func (c *GitHubClient) DeleteRepo(repoName string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repoName)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending delete request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error deleting repository: %s, response: %s", resp.Status, string(body))
	}

	return nil
}

func main() {
	githubToken := os.Getenv("GHUB_TOKEN")
	if githubToken == "" {
		log.Fatal("Error: GHUB_TOKEN environment variable is not set")
	}

	if len(os.Args) < 2 {
		log.Fatal("Error: Repository name argument is required")
	}

	repoName := os.Args[1]
	client := &GitHubClient{
		HTTPClient: &http.Client{},
		Token:      githubToken,
	}

	err := client.DeleteRepo(repoName)
	if err != nil {
		log.Fatalf("Failed to delete repository: %v", err)
	}

	fmt.Println("Repository deleted successfully:", repoName)
}
