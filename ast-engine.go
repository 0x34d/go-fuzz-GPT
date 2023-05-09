package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
)

func getHarnessFunction(files []*ast.File) []*ast.FuncDecl {
	var harnessFuncs []*ast.FuncDecl

	for _, file := range files {
		// Iterate over all function declarations in the file
		for _, decl := range file.Decls {
			// Check if the declaration is a function declaration
			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			// Check if the function has string, []byte, or io.Reader arguments
			for _, field := range funcDecl.Type.Params.List {
				switch fieldType := field.Type.(type) {
				case *ast.Ident:
					if fieldType.Name == "string" {
						harnessFuncs = append(harnessFuncs, funcDecl)
						break
					}
				case *ast.ArrayType:
					if arrayType, ok := fieldType.Elt.(*ast.Ident); ok && arrayType.Name == "byte" {
						harnessFuncs = append(harnessFuncs, funcDecl)
						break
					}
				case *ast.SelectorExpr:
					if xIdent, ok := fieldType.X.(*ast.Ident); ok && xIdent.Name == "io" && fieldType.Sel.Name == "Reader" {
						harnessFuncs = append(harnessFuncs, funcDecl)
						break
					}
				}
			}
		}
	}

	return harnessFuncs
}

func analysis(goAstFiles []*ast.File, testAstFiles []*ast.File) {
	goHarnessFuncs := getHarnessFunction(goAstFiles)

	// Print the names and receivers of the harness functions
	if len(goHarnessFuncs) > 0 {
		funcMap := make(map[string][]*ast.FuncDecl)

		// Group similar function names with different signatures
		for _, funcDecl := range goHarnessFuncs {
			funcName := funcDecl.Name.Name
			funcMap[funcName] = append(funcMap[funcName], funcDecl)
		}

		// Collect functions with tests as strings
		for funcName, funcDecls := range funcMap {
			count, callDetails := countFunctionCalls(funcName, testAstFiles)

			// Process only functions with at least one test
			if count > 0 {
				var functions, tests string
				for _, funcDecl := range funcDecls {
					functions += "\nFunction:\n" + nodeString(funcDecl) + "\n"
				}
				tests = "\nTest functions:" + callDetails
				gptWork(funcName, functions, tests)
			}
		}
	}
}

func countFunctionCalls(funcName string, testAstFiles []*ast.File) (int, string) {
	counter := 0
	callDetails := ""

	for _, testFile := range testAstFiles {
		var funcStack []*ast.FuncDecl

		ast.Inspect(testFile, func(n ast.Node) bool {
			if f, ok := n.(*ast.FuncDecl); ok {
				// Push the function onto the stack when we enter its scope
				funcStack = append(funcStack, f)
			}

			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Check if the function call is an ast.Ident or ast.SelectorExpr
			switch fn := call.Fun.(type) {
			case *ast.Ident:
				if fn.Name == funcName {
					counter++
					if len(funcStack) > 0 {
						// Get the enclosing function from the top of the stack
						enclosingFunc := funcStack[len(funcStack)-1]
						callDetails += fmt.Sprintf("\n%s\n", nodeString(enclosingFunc))
					}
				}
			case *ast.SelectorExpr:
				if fn.Sel.Name == funcName {
					counter++
					if len(funcStack) > 0 {
						// Get the enclosing function from the top of the stack
						enclosingFunc := funcStack[len(funcStack)-1]
						callDetails += fmt.Sprintf("\n%s\n", nodeString(enclosingFunc))
					}
				}
			}

			// Pop the function from the stack when we exit its scope
			if _, ok := n.(*ast.FuncDecl); ok && len(funcStack) > 0 {
				funcStack = funcStack[:len(funcStack)-1]
			}

			return true
		})
	}
	return counter, callDetails
}

// Helper function to stringify an AST node
func nodeString(node ast.Node) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), node)
	if err != nil {
		return ""
	}
	return buf.String()
}
