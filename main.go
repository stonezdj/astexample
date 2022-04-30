package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"

	// service
	"io/ioutil"
	"path/filepath"
)

func main() {

	// prepare

	////////////////////////////////////////////////
	//
	dir := createDir()
	defer os.RemoveAll(dir) // clean up
	//
	////////////////////////////////////////////////

	// go

	fs := token.NewFileSet()
	pkgs, err := parser.ParseDir(fs, dir, nil, 0)
	if err != nil {
		log.Fatal("parsing dir:", err)
	}

	for name, pkg := range pkgs {
		fmt.Println("Found package:", name)
		ast.Walk(VisitorFunc(FindTypes), pkg)
	}
}

type VisitorFunc func(n ast.Node) ast.Visitor

func (f VisitorFunc) Visit(n ast.Node) ast.Visitor {
	return f(n)
}

func FindTypes(n ast.Node) ast.Visitor {
	switch n := n.(type) {
	case *ast.Package:
		return VisitorFunc(FindTypes)
	case *ast.File:
		return VisitorFunc(FindTypes)
	case *ast.GenDecl:
		if n.Tok == token.TYPE {
			return VisitorFunc(FindTypes)
		}
	case *ast.TypeSpec:
		fmt.Println(n.Name.Name)
	}
	return nil
}

// ======================================//
// create tmp dir with following content //
// ======================================//

const (
	someGo = `package some

import (
	"fmt"
)

type String string

type Text String

type Age uint

type Some struct {
	Name   Text
	Age    Age
	hidden int
	Any    Any
	fmt.Stringer
}

`
	anyGo = `package some

type ApplyFunc func(string) error

type Weight int

type Any struct {
	WeightChan chan Weight
	Apply ApplyFunc
	Slice []Age
}

`
)

func createFile(dir, name, content string) {
	tmpfn := filepath.Join(dir, name)
	if err := ioutil.WriteFile(tmpfn, []byte(content), 0644); err != nil {
		log.Fatal(err)
	}
}

func createDir() string {
	dir, err := ioutil.TempDir("", "some")
	if err != nil {
		log.Fatal(err)
	}

	createFile(dir, "some.go", someGo)
	createFile(dir, "any.go", anyGo)
	return dir
}
