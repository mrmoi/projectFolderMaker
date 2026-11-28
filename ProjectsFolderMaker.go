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

// pickFile asks the user to choose the input text file. On macOS it shows the
// native file dialog; if that is unavailable or cancelled, or on other systems,
// it falls back to typing a path at the terminal (Enter keeps defaultPath).
func pickFile(defaultPath string) (string, error) {
	if runtime.GOOS == "darwin" {
		out, err := exec.Command("osascript", "-e",
			`POSIX path of (choose file with prompt "Select the project list text file")`).Output()
		if err == nil {
			if path := strings.TrimSpace(string(out)); path != "" {
				return path, nil
			}
		}
	}

	fmt.Printf("Path to text file [%s]: ", defaultPath)
	line, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil && line == "" {
		return "", fmt.Errorf("no file selected: %w", err)
	}
	line = strings.Trim(strings.TrimSpace(line), `"'`)
	if line == "" {
		return defaultPath, nil
	}
	return line, nil
}

func main() {
	inputPath := flag.String("input", "/Users/Moises.Martinez/Downloads/q4ProjectList.txt", "text file containing folder names, one per line")
	baseDir := flag.String("base", "/Users/Moises.Martinez/Downloads/Projects/q4-2026/", "base directory where folders are created")
	pick := flag.Bool("pick", false, "choose the input text file from a prompt instead of using -input")
	flag.Parse()

	if *pick {
		chosen, err := pickFile(*inputPath)
		if err != nil {
			log.Fatalf("Failed to pick file: %v", err)
		}
		*inputPath = chosen
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
