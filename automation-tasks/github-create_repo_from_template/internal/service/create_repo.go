package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	apiURL = "https://api.github.com"
)

func CreateRepoFromTempl(template string, repo_name string) error {
	url := fmt.Sprintf("%s/repos/%s/generate", apiURL, template)
	log.Println("URL ", url)

	ghub_token := os.Getenv("GHUB_TOKEN")
	if ghub_token == "" {
		log.Println("Error: GHUB_TOKEN not set")
		return errors.New("GHUB_TOKEN not set")
	}

	type repoDetails struct {
		Name               string `json:"name"`
		Description        string `json:"description"`
		Private            bool   `json:"private"`
		IncludeAllBranches bool   `json:"include_all_branches"`
	}

	repodDetails := repoDetails{
		Name:               repo_name,
		Description:        "This is a test java repo",
		Private:            false,
		IncludeAllBranches: true,
	}

	repoJson, err := json.Marshal(repodDetails)
	if err != nil {
		log.Println("Error marshalling repo details: ", err)
		return fmt.Errorf("error marshalling repo details: %w", err)

	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(repoJson))
	if err != nil {
		log.Println("Error creating request: ", err)
		return fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+ghub_token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	httpClient := &http.Client{}

	resp, err := httpClient.Do(req)

	if err != nil {
		log.Println("Error sending request: ", err)
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		log.Println("Error creating repository: ", resp)
		return fmt.Errorf("error creating repository: %s", resp.Status)
	}
	log.Println("Repository created successfully")

	return nil

}
