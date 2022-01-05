package main

import (
	"bytes"
	"fmt"
	"github.com/tmr232/go-explore/itertools"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"log"
	"strings"
	"testing"
	"text/template"
)

func _TestGotoSyntax(t *testing.T) {
	fset := token.NewFileSet()
	src := `
package src

func f(flag bool) bool {
	if flag {
		goto retFalse
	}
	return true

retFalse:
	return false

}
`

	af, err := parser.ParseFile(fset, "test.go", src, 0)
	if err != nil {
		log.Fatal(err)
	}

	print := func() {
		var out bytes.Buffer
		format.Node(&out, fset, af)
		fmt.Println(out.String())
	}

	print()
	ast.Print(fset, af)
}

func TestAst(t *testing.T) {
	fset := token.NewFileSet()
	src := `
package src

func f() {
	x := 1
	if x == 2 {
		fmt.Println("Two!")
	} else if x == 1 {
		fmt.Println("One!")
	} else {
		fmt.Println("Unknown!")
	}
	fmt.Println("We're done!")
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
	flt := simpleFlattener{fset: fset}
	ast.Walk(visitFunc(flt.printFlattenedIfs), af)
}

func invertConditions(n ast.Node) ast.Visitor {
	is, ok := n.(*ast.IfStmt)
	if ok {
		invertCondition(is)
	}
	return visitFunc(invertConditions)
}

func invertCondition(is *ast.IfStmt) {
	cond := is.Cond
	is.Cond = &ast.UnaryExpr{Op: token.NOT, X: cond}
}

func flattenIfs(c *astutil.Cursor) bool {
	block, ok := c.Node().(*ast.BlockStmt)
	if ok {
		for i, stmt := range block.List {
			_, ok := stmt.(*ast.IfStmt)
			if !ok {
				continue
			}
			newBlock := &ast.BlockStmt{List: []ast.Stmt{stmt}}
			block.List[i] = newBlock
		}
	}
	return true
}

/*
The idea is to just build a new AST and ignore the original.
We don't need to replace anything anyhow, as we're generating a new file.
*/
type simpleFlattener struct {
	fset *token.FileSet
	id   int
}

type ifData struct {
	Cond      string
	TrueName  string
	TrueBody  string
	FalseName string
	FalseBody string
	EndName   string
	Nop       string
}

func renderNode(fset *token.FileSet, node ast.Node) string {
	var out bytes.Buffer
	format.Node(&out, fset, node)
	return out.String()
}

func (flt *simpleFlattener) makeId(name string) string {
	flt.id++
	return fmt.Sprintf("__%s_%d", name, flt.id)
}

func (flt *simpleFlattener) flatten(stmt ast.Node) string {
	switch stmt.(type) {
	case *ast.IfStmt:
		return flt.flattenIf(stmt.(*ast.IfStmt))
	case *ast.BlockStmt:
		builder := strings.Builder{}
		for _, s := range stmt.(*ast.BlockStmt).List {
			builder.WriteString(flt.flatten(s))
			builder.WriteString("\n")
		}
		return builder.String()
	//case *ast.ReturnStmt:
	//	builder := strings.Builder{}
	//	next := flt.makeId("Next")
	//	builder.WriteString(fmt.Sprintf("__next = "))
	default:
		return renderNode(flt.fset, stmt)
	}
}

func (flt *simpleFlattener) flattenIf(ifStmt *ast.IfStmt) string {

	data := ifData{
		Cond:      renderNode(flt.fset, ifStmt.Cond),
		TrueName:  flt.makeId("onTrue"),
		TrueBody:  flt.flatten(ifStmt.Body),
		FalseName: flt.makeId("onFalse"),
		FalseBody: flt.flatten(ifStmt.Else),
		EndName:   flt.makeId("end"),
		Nop:       "nop()",
	}

	tmpl, err := template.New("If").Parse(
		`if {{.Cond}} {
	goto {{.TrueName}}
} else {
	goto {{.FalseName}}
}
{{.TrueName}}:
{{.TrueBody}}
goto {{.EndName}}
{{.FalseName}}:
{{.FalseBody}}
goto {{.EndName}}
{{.EndName}}:
{{.Nop}}
`,
	)
	if err != nil {
		panic(err)
	}

	var out bytes.Buffer
	err = tmpl.Execute(&out, data)
	if err != nil {
		panic(err)
	}
	return out.String()
}
func (flt *simpleFlattener) printFlattenedIfs(n ast.Node) ast.Visitor {
	fd, ok := n.(*ast.FuncDecl)
	if ok {
		fmt.Println("====================================")
		fmt.Println(flt.flatten(fd.Body))
	}
	return visitFunc(flt.printFlattenedIfs)
}

func MyGen() itertools.Iterator[int] {
	__next := 0
	var __zero int
	__nop := func() {}
	advance := func() (bool, int) {
		switch __next {

		case 0:
			goto __next_0

		case 1:
			goto __next_1

		case 2:
			goto __next_2

		case 3:
			goto __next_3

		}

	__next_0:

		__next = 1
		return true, 1
	__next_1:
		__nop()

		__next = 2
		return true, 2
	__next_2:
		__nop()

		__next = 3
		return true, 3
	__next_3:
		__nop()

		return false, __zero
	}
	return itertools.FromAdvance(advance)
}

func TestMyGen(t *testing.T) {
	for gen := MyGen(); gen.Next(); {
		fmt.Println(gen.Value())
	}
}
