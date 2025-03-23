package utils

import (
	"log"
	"os"
	"strings"
)

func GetFileContent(filePath string) []string {
	file, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("unable to read file, err: %v", err)
		return nil
	}
	content := strings.Split(string(file), "\n")
	return content
}

func LoadContentFromFile(filepath string) []string {
	content := GetFileContent(filepath)
	if content == nil {
		log.Fatalf("unable to load content from file: %s", filepath)
	}
	return content
}
