package main

import (
	"flag"
	"fmt"
	"github-create_branch/internal/service"
)

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
	err := service.CreateBranch(*repoName, *baseBranch, *newBranch)
	if err != nil {
		fmt.Println("Error creating branch:", err)
	}
}