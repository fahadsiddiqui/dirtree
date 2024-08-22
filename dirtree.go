package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	gitignore "github.com/sabhiram/go-gitignore"
)

func main() {
	var rootDir string

	switch len(os.Args) {
	case 1:
		rootDir = "."
	case 2:
		rootDir = strings.Join(os.Args[1:2], "")
		_, err := os.Stat(rootDir)
		if os.IsNotExist(err) {
			fmt.Println("File does not exist")
			os.Exit(1)
		}
		if err != nil {
			fmt.Printf("Error %v", err)
		}
	default:
		fmt.Printf("Usage: dirtree [directory]")
		os.Exit(1)
	}

	fmt.Println(".")
	printTree(rootDir, "", nil)
}

func shouldIgnore(path string, ignoreMatchers []*gitignore.GitIgnore) bool {
	for _, matcher := range ignoreMatchers {
		if matcher.MatchesPath(path) {
			return true
		}
	}
	return false
}

func loadIgnoreMatchers(path string, parentMatchers []*gitignore.GitIgnore) []*gitignore.GitIgnore {
	ignoreFile := filepath.Join(path, ".gitignore")
	if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
		return parentMatchers
	}

	ignoreMatcher, err := gitignore.CompileIgnoreFile(ignoreFile)
	if err != nil {
		fmt.Printf("Error reading .gitignore in %s: %v\n", path, err)
		return parentMatchers
	}

	return append(parentMatchers, ignoreMatcher)
}

func printTree(root string, prefix string, parentMatchers []*gitignore.GitIgnore) {
	ignoreMatchers := loadIgnoreMatchers(root, parentMatchers)

	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", root, err)
		return
	}

	for i, file := range files {
		if file.Name() == ".git" {
			continue
		}

		relPath, _ := filepath.Rel(".", filepath.Join(root, file.Name()))
		if shouldIgnore(relPath, ignoreMatchers) {
			continue
		}

		if i == len(files)-1 {
			fmt.Printf("%s└── %s\n", prefix, file.Name())
		} else {
			fmt.Printf("%s├── %s\n", prefix, file.Name())
		}

		if file.IsDir() {
			newPrefix := prefix
			if i == len(files)-1 {
				newPrefix += "    "
			} else {
				newPrefix += "│   "
			}
			printTree(filepath.Join(root, file.Name()), newPrefix, ignoreMatchers)
		}
	}
}
