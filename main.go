package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	Blue         = "\033[34m"
	CyanColor    = "\033[36m"
	GreenColor   = "\033[32m"
	MagentaColor = "\033[35m"
	Red          = "\033[31m"
	ResetColor   = "\033[0m"
	White        = "\033[37m"
	YellowColor  = "\033[33m"
)

var RemoteGitURL string

func getRemoteOriginURL(repoPath string) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		RemoteGitURL = "This project does not have any public repo"
		return
	}

	url := strings.TrimSpace(string(output))
	RemoteGitURL = url
}

func parseFiles(files []string) ([]*ast.File, error) {
	fset := token.NewFileSet()

	var astFiles []*ast.File

	for _, file := range files {
		astFile, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("error parsing file %s: %v", file, err)
		}

		astFiles = append(astFiles, astFile)
	}

	return astFiles, nil
}

func separateFiles(files []string) (goFiles []string, testFiles []string) {
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			testFiles = append(testFiles, file)
		} else {
			goFiles = append(goFiles, file)
		}
	}

	return goFiles, testFiles
}

func processFiles(dir string, files []string) {
	goFiles, testFiles := separateFiles(files)

	fmt.Printf(GreenColor+"\nDirectory: %s\n"+ResetColor, dir)

	if len(goFiles) > 0 {
		fmt.Println(YellowColor + "Go Files:" + ResetColor)
		for _, file := range goFiles {
			fmt.Printf(YellowColor+"\t%s\n"+ResetColor, file)
		}
	}

	if len(testFiles) > 0 {
		fmt.Println(CyanColor + "Test Files:" + ResetColor)
		for _, file := range testFiles {
			fmt.Printf(CyanColor+"\t%s\n"+ResetColor, file)
		}
		fmt.Println()
	}

	goAstFiles, goErr := parseFiles(goFiles)
	testAstFiles, testErr := parseFiles(testFiles)

	if goErr != nil {
		fmt.Printf("Error parsing go files: %v\n", goErr)
		return
	}

	if testErr != nil {
		fmt.Printf("Error parsing test files: %v\n", testErr)
		return
	}

	analysis(goAstFiles, testAstFiles)
}

func processDirectory(path string) {
	dirEntry, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("Error reading directory: %v\n", err)
		return
	}

	var files []string
	for _, entry := range dirEntry {
		entryPath := filepath.Join(path, entry.Name())
		if entry.IsDir() {
			processDirectory(entryPath)
		} else {
			if strings.HasSuffix(entry.Name(), ".go") {
				files = append(files, entryPath)
			}
		}
	}

	if len(files) > 0 {
		getRemoteOriginURL(path)
		processFiles(path, files)
	}
}

func main() {
	args := os.Args

	if len(args) != 2 {
		fmt.Println("Provide input as <Directory>")
		return
	}

	path := args[1]
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	if !fileInfo.IsDir() {
		fmt.Printf("The provided path is not a directory: %s\n", path)
		return
	}

	fmt.Printf("Processing directory: %s\n", path)
	processDirectory(path)
}
