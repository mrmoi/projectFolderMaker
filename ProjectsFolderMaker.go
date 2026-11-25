package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	// 1. Open the text file containing the folder names
	file, err := os.Open("/Users/Moises.Martinez/Downloads/q4ProjectList.txt")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// Optional: Specify a base directory where folders should be created
	baseDir := "/Users/Moises.Martinez/Downloads/Projects/q4-2026/"

	// 2. Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		var quoRe = regexp.MustCompile(`\bQUO-\d{5}-[A-Z0-9]{4}\b`)
		folderName := quoRe.ReplaceAllString(scanner.Text(), "")
		if folderName == "" {
			continue
		}

		//Trim whitespace and skip empty lines
		folderNameTrimmed := strings.TrimSpace(scanner.Text())
		if folderNameTrimmed == "" {
			continue
		}

		// Clean and join the path safely for any OS
		targetPath := filepath.Join(baseDir, filepath.Clean(folderName))

		// 3. Create the folder (and any nested parent folders)
		err := os.MkdirAll(targetPath, 0755)
		if err != nil {
			log.Printf("Error creating folder %q: %v", targetPath, err)
			continue
		}

		fmt.Printf("Created: %s\n", targetPath)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
}
