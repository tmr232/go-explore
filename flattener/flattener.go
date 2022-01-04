package main

import (
	"bytes"
	"fmt"
	"github.com/edwingeng/deque"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"log"
	"strings"
)

type visitFunc func(ast.Node) ast.Visitor

func (f visitFunc) Visit(n ast.Node) ast.Visitor { return f(n) }

type Flattener struct {
	fset    *token.FileSet
	stateId int
}

func (flt *Flattener) render(node ast.Node) string {
	var out bytes.Buffer
	format.Node(&out, flt.fset, node)
	return out.String()
}

func (flt *Flattener) getStateId() int {
	flt.stateId++
	return flt.stateId
}

func (flt *Flattener) setNext() (setter, label string) {
	stateId := flt.getStateId()
	label = fmt.Sprintf("__next_%d:", stateId)
	setter = fmt.Sprintf("__next = %d", stateId)
	return
}

func (flt *Flattener) generateStateSwitch() string {
	builder := strings.Builder{}

	builder.WriteString("switch __next {\n")

	for i := 0; i <= flt.stateId; i++ {
		builder.WriteString(fmt.Sprintf("case %[1]d:\n    goto __next_%[1]d\n", i))
	}

	builder.WriteString("}")

	return builder.String()
}

func (flt *Flattener) flatten(node ast.Node) string {
	switch node.(type) {
	case *ast.ReturnStmt:
		builder := strings.Builder{}
		setter, label := flt.setNext()
		builder.WriteString(setter)
		builder.WriteString("\n")
		builder.WriteString(flt.render(node))
		builder.WriteString("\n")
		builder.WriteString(label)
		builder.WriteString("\n__nop()\n")
		return builder.String()
	case *ast.BlockStmt:
		builder := strings.Builder{}
		for _, stmt := range node.(*ast.BlockStmt).List {
			builder.WriteString(flt.flatten(stmt))
		}
		return builder.String()
	}
	return "// UNSUPPORTED\n"
}

func (flt *Flattener) FlattenFunction(fd *ast.FuncDecl) string {
	builder := strings.Builder{}

	builder.WriteString(flt.render(fd.Type))
	builder.WriteString("{\n")
	builder.WriteString("__nop = func() {}\n")
	builder.WriteString("iterator = func() Iterator[int]")
	builder.WriteString(flt.flatten(fd.Body))
	builder.WriteString("}\n")

	return builder.String()
}

func (flt *Flattener) printFlattenedIfs(n ast.Node) ast.Visitor {
	fd, ok := n.(*ast.FuncDecl)
	if ok {
		fmt.Println("====================================")
		fmt.Println(flt.FlattenFunction(fd))
		return nil
	}
	return visitFunc(flt.printFlattenedIfs)
}

func collectFuncDecls(node ast.Node, recurse bool) []*ast.FuncDecl {
	dq := deque.NewDeque()
	var visitor func(ast.Node) ast.Visitor
	visitor = func(n ast.Node) ast.Visitor {
		_, ok := n.(*ast.FuncDecl)
		if ok {
			dq.PushBack(n)
			if !recurse {
				return nil
			}
		}
		return visitFunc(visitor)
	}

	ast.Walk(visitFunc(visitor), node)

	if dq.Empty() {
		return nil
	}

	result := make([]*ast.FuncDecl, dq.Len())
	for i, elem := range dq.DequeueMany(0) {
		result[i] = elem.(*ast.FuncDecl)
	}
	return result
}

func main() {
	fset := token.NewFileSet()
	src := `
package src

func f() int {
	return 1
	return 2
	return 3
}
`

	af, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		log.Fatal(err)
	}
	//
	//print := func() {
	//	var out bytes.Buffer
	//	format.Node(&out, fset, af)
	//	fmt.Println(out.String())
	//}
	//
	//print()
	//
	//ast.Walk(visitFunc(invertConditions), af)
	//astutil.Apply(af, nil, flattenIfs)
	//print()

	funcDecls := collectFuncDecls(af, false)

	flt := Flattener{fset: fset}
	fmt.Println(flt.FlattenFunction(funcDecls[0]))
}
