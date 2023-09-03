package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
)

// Get harness functions
func GetHarnessFunction(codefiles []*ast.File) []*ast.FuncDecl {
	var harnessFuncs []*ast.FuncDecl

	for _, codefile := range codefiles {
		for _, decl := range codefile.Decls {
			funcDecl, ok := decl.(*ast.FuncDecl) // check if `decl` is a `ast.FuncDecl`.
			if !ok {
				continue
			}

			for _, field := range funcDecl.Type.Params.List {
				switch fieldtype := field.Type.(type) {
				case *ast.Ident:
					if fieldtype.Name == "string" {
						harnessFuncs = append(harnessFuncs, funcDecl)
						break
					}
				case *ast.ArrayType:
					if arrayType, ok := fieldtype.Elt.(*ast.Ident); ok && arrayType.Name == "byte" {
						harnessFuncs = append(harnessFuncs, funcDecl)
						break
					}
				case *ast.SelectorExpr:
					if xIdent, ok := fieldtype.X.(*ast.Ident); ok && xIdent.Name == "io" && fieldtype.Sel.Name == "Reader" {
						harnessFuncs = append(harnessFuncs, funcDecl)
						break
					}
				}
			}
		}
	}

	return harnessFuncs
}

// Get test function corresponding to the function name
func countFunctionCalls(funcname string, testfiles []*ast.File) (int, string) {
	var counter int
	var callDetails string
	var funcstack []*ast.FuncDecl

	for _, testfile := range testfiles {
		// Inspect the AST of the test file
		ast.Inspect(testfile, func(n ast.Node) bool {
			// Append function declaration to the function stack
			if f, ok := n.(*ast.FuncDecl); ok {
				funcstack = append(funcstack, f)
			}

			// Check if the node is a function call expression
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			// Not more than 3 test functions
			if counter == 3 {
				return true
			}

			switch fn := call.Fun.(type) {
			// Standalone function call
			case *ast.Ident:
				if fn.Name == funcname {
					counter++
					if len(funcstack) > 0 {
						enclosingFunc := funcstack[len(funcstack)-1]
						callDetails += fmt.Sprintf("\n%s\n", nodeString(enclosingFunc))
					}
				}
			// Method function call
			case *ast.SelectorExpr:
				if fn.Sel.Name == funcname {
					counter++
					if len(funcstack) > 0 {
						enclosingFunc := funcstack[len(funcstack)-1]
						callDetails += fmt.Sprintf("\n%s\n", nodeString(enclosingFunc))
					}
				}
			}

			// Pop the function declaration from the function stack
			if _, ok := n.(*ast.FuncDecl); ok && len(funcstack) > 0 {
				funcstack = funcstack[:len(funcstack)-1]
			}

			return true
		})
	}

	return counter, callDetails
}

// Get the string representation of the node
func nodeString(node ast.Node) string {
	var buf bytes.Buffer
	err := printer.Fprint(&buf, token.NewFileSet(), node)
	if err != nil {
		return ""
	}
	return buf.String()
}

// Analysis the AST
func Analysis(codeast []*ast.File, testast []*ast.File) {
	harnessFuncs := GetHarnessFunction(codeast)

	if len(harnessFuncs) > 0 {
		funcMap := make(map[string][]*ast.FuncDecl)

		for _, funcDecl := range harnessFuncs {
			funcName := funcDecl.Name.Name
			funcMap[funcName] = append(funcMap[funcName], funcDecl)
		}

		for funcName, funcDecls := range funcMap {
			count, callDetails := countFunctionCalls(funcName, testast)

			if count > 0 {
				var functions, tests string
				for _, funcDecl := range funcDecls {
					functions += "\nFunction:\n" + nodeString(funcDecl) + "\n"
				}
				tests = "\nTest functions:" + callDetails
				GPTWork(funcName, functions, tests)
			}
		}
	}
}
