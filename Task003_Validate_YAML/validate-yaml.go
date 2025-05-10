package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: validate-yaml <file.yaml>")
		os.Exit(1)
	}
	filePath := os.Args[1]
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("Error in reading file:", err)
		os.Exit(1)
	}
	var parsedData map[string]interface{}
	err = yaml.Unmarshal(data, &parsedData)
	if err != nil {
		fmt.Println("Error in parsing YAML:", err)
		os.Exit(1)
	}
	fmt.Println("YAML syntax is valid")
}
