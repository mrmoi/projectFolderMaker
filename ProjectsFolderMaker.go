package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var quoRe = regexp.MustCompile(`\bQUO-\d{5}-[A-Z0-9]{4}\b`)

// openFile opens path in the system's default application.
func openFile(path string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", path)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", "", path)
	default:
		cmd = exec.Command("xdg-open", path)
	}
	return cmd.Start()
}

func main() {
	inputPath := flag.String("input", "/Users/Moises.Martinez/Downloads/q4ProjectList.txt", "text file containing folder names, one per line")
	baseDir := flag.String("base", "/Users/Moises.Martinez/Downloads/Projects/q4-2026/", "base directory where folders are created")
	openOnly := flag.Bool("open", false, "open the input text file in the default application and exit")
	flag.Parse()

	if *openOnly {
		if _, err := os.Stat(*inputPath); err != nil {
			log.Fatalf("Cannot open file: %v", err)
		}
		if err := openFile(*inputPath); err != nil {
			log.Fatalf("Failed to open %q: %v", *inputPath, err)
		}
		fmt.Printf("Opened: %s\n", *inputPath)
		return
	}

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
		folderName = strings.NewReplacer("/", "-", `\`, "-").Replace(folderName)
		if folderName == "." || folderName == ".." {
			log.Printf("Skipping invalid folder name %q", folderName)
			continue
		}

		// Replace hyphen and double spaces
		folderName = strings.NewReplacer("-", "", "  ", "").Replace(folderName)
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
