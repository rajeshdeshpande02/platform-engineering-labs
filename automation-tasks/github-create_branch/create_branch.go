package main

import (
	"flag"
	"fmt"
	"os"
	"net/http"
	"bytes"
	"encoding/json"
	
)

const (
	apiURL = "https://api.github.com"
	owner = "rajeshdeshpande02"
)

type Result struct {
	Object struct {
		Sha string `json:"sha"`
	} `json:"object"`
}

func main() {
	repoName := flag.String("repo-name", "", "The name of the repository in owner/repo format")
	baseBranch := flag.String("base-branch", "", "The base branch to create the new branch from")
	newBranch := flag.String("new-branch", "", "The name of the new branch to create")

	flag.Parse()

	//Quick and dirty validation. Could be better
	if *baseBranch == "" || *newBranch == "" || *repoName == "" {
		fmt.Println("Error: repo-name, base-branch and new-branch flags are required")
		return
	}

	fmt.Printf("Creating new branch '%s' from base branch '%s'\n", *newBranch, *baseBranch)

	ghub_token := os.Getenv("GHUB_TOKEN")
	if ghub_token == "" {
		fmt.Println("Error getting GHUB_TOKEN from environment variables")
		return
	}
	baseBranchSHA, err := getBaseBranchSHA(*repoName, *baseBranch)
	if err != nil {
		fmt.Println("Error getting base branch SHA:", err)
		return
	}

	err = createBranch(*repoName, *newBranch, baseBranchSHA)
	if err != nil {
		fmt.Println("Error creating branch:", err)
		return
	}

}

func getBaseBranchSHA(repoName string, baseBranch string) (string, error) {
	url := fmt.Sprintf("%s/repos/%s/git/refs/heads/%s", apiURL, repoName, baseBranch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

    httpClient := &http.Client{}
	
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error getting base branch SHA: %s", resp.Status)
	}

	
	var result Result
	
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}
   
	fmt.Println("Base branch SHA:", result.Object.Sha)

	return result.Object.Sha, nil

}


func createBranch(repoName string, newBranch string, baseBranchSHA string) error {
	url := fmt.Sprintf("%s/repos/%s/git/refs", apiURL, repoName)

	body := map[string]string{
		"ref":  "refs/heads/" + newBranch,
		"sha":  baseBranchSHA,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("GHUB_TOKEN"))
	req.Header.Set("content-type", "application/json")

	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("error creating branch: %s", resp.Status)
	}
	fmt.Println("Branch created successfully:", newBranch)
	return nil
}