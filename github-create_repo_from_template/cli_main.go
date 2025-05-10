package main

import (
	"flag"
	"github-create_repo_from_template/internal/service"
	"log"
)

func main() {
	template := flag.String("template", "", "The template repository to use in the format <owner>/<repo>")
	repo_name := flag.String("repo_name", "", "The name of the new repository")
	flag.Parse()
	if *template == "" || *repo_name == "" {
		flag.Usage()
		return
	}
	err := service.CreateRepoFromTempl(*template, *repo_name)
	if err != nil {
		log.Fatalf("Error creating repository from template: %v", err)
	}

}
