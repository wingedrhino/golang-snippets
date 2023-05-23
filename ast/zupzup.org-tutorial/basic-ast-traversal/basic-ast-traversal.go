package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"

	"github.com/davecgh/go-spew/spew"
)

const (
	filename = "/home/gautham/ext/workspace/golang-snippets/ast/zupzup.org-tutorial/basic-ast-traversal/sample-code/sample-code.go"
)

func main() {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	spew.Printf("\n\n\nFile contains: %v\n\n\n\n", file)

	for _, i := range file.Imports {
		fmt.Printf("Import found: %s\n", i.Path.Value)
	}

	for _, d := range file.Decls {
		switch v := d.(type) {
		case *ast.FuncDecl:
			fmt.Printf("Found function with name %s\n", v.Name.Name)
			rcv := v.Recv
			if rcv == nil {
				fmt.Printf("Simple function found, not a method.\n")
			} else {
				fields := rcv.List
				fmt.Printf("Method has %d receivers\n", len(fields))
				for _, f := range fields {
					names := make([]string, len(f.Names))
					for _, name := range f.Names {
						names = append(names, name.Name)
					}
				}
			}
		}
	}

}
