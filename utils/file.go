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
