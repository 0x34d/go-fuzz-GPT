package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// Constants for color output
const (
	ResetColor = "\033[0m"
	RedColor   = "\033[31m"
	BlueColor  = "\033[34m"
	CyanColor  = "\033[36m"
	GreenColor = "\033[32m"
)

// Parse the files to get AST
func ParseFiles(files []string) ([]*ast.File, error) {
	var astfiles []*ast.File
	fset := token.NewFileSet()

	for _, file := range files {
		astfile, err := parser.ParseFile(fset, file, nil, 0)
		if err != nil {
			return nil, fmt.Errorf("error: %s: %v", file, err)
		}

		astfiles = append(astfiles, astfile)
	}

	return astfiles, nil
}

// Separate code.go and _test.go files
func SeparateFiles(files []string) (codefiles []string, testfiles []string) {
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			testfiles = append(testfiles, file)
		} else {
			codefiles = append(codefiles, file)
		}
	}

	return codefiles, testfiles
}

// Process files and send for analysis
func ProcessFiles(files []string) {
	codefiles, testfiles := SeparateFiles(files)
	if len(codefiles) == 0 || len(testfiles) == 0 {
		return
	}

	codeast, codeerr := ParseFiles(codefiles)
	testast, testerr := ParseFiles(testfiles)
	if codeerr != nil || testerr != nil {
		return
	}

	Analysis(codeast, testast)
}

// Process Dir structure
func ProcessDir(path string) {
	var files []string

	direntry, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf(RedColor+"Error reading directory: %v\n"+ResetColor, err)
		return
	}

	for _, entry := range direntry {
		entrypath := filepath.Join(path, entry.Name())

		if entry.IsDir() {
			ProcessDir(entrypath)
			continue
		}

		if strings.HasSuffix(entry.Name(), ".go") {
			files = append(files, entrypath)
		}
	}

	if len(files) > 0 {
		fmt.Printf(GreenColor+"\nDirectory: %s\n"+ResetColor, path)
		ProcessFiles(files)
	}
}

// Main function
func main() {
	if len(os.Args) != 2 {
		fmt.Printf(RedColor + "Provide input as <dir> \n" + ResetColor)
		return
	}
	ProcessDir(os.Args[1])
}
