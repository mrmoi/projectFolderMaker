package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var quoRe = regexp.MustCompile(`\bQUO-\d{5}-[A-Z0-9]{4}\b`)

func main() {
	inputPath := flag.String("input", "/Users/Moises.Martinez/Downloads/q4ProjectList.txt", "text file containing folder names, one per line")
	baseDir := flag.String("base", "/Users/Moises.Martinez/Downloads/Projects/q4-2026/", "base directory where folders are created")
	flag.Parse()

	// 1. Open the text file containing the folder names
	file, err := os.Open(*inputPath)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	// 2. Read the file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Remove quote IDs, then trim whitespace and skip empty names
		folderName := strings.TrimSpace(quoRe.ReplaceAllString(scanner.Text(), ""))
		if folderName == "" {
			continue
		}

		// Neutralize path separators so a name can't escape baseDir
		folderName = strings.NewReplacer("-", "", "-", "", "  ", "").Replace(folderName)
		if folderName == "." || folderName == ".." {
			log.Printf("Skipping invalid folder name %q", folderName)
			continue
		}

		targetPath := filepath.Join(*baseDir, folderName)

		// 3. Create the folder (and any missing parent folders)
		if err := os.MkdirAll(targetPath, 0755); err != nil {
			log.Printf("Error creating folder %q: %v", targetPath, err)
			continue
		}

		fmt.Printf("Created: %s\n", targetPath)
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
}
