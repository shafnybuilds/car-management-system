package driver

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestInitDBDoesNotShadowDatabaseHandle(t *testing.T) {
	file, err := parser.ParseFile(token.NewFileSet(), "postgress.go", nil, 0)
	if err != nil {
		t.Fatal(err)
	}
	globalDB := file.Scope.Lookup("db")
	if globalDB == nil {
		t.Fatal("package database handle is missing")
	}
	for _, decl := range file.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Name.Name != "InitDB" {
			continue
		}
		ast.Inspect(fn.Body, func(node ast.Node) bool {
			if id, ok := node.(*ast.Ident); ok && id.Name == "db" && id.Obj != globalDB {
				t.Errorf("InitDB shadows the package database handle at offset %d", id.Pos())
			}
			return true
		})
	}
}
