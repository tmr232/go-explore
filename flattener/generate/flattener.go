package main

import (
	"bytes"
	"fmt"
	"github.com/edwingeng/deque"
	"github.com/tmr232/go-explore/itertools"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/packages"
	"log"
	"os"
	"strings"
)

type visitFunc func(ast.Node) ast.Visitor

func (f visitFunc) Visit(n ast.Node) ast.Visitor { return f(n) }

type Flattener struct {
	fset      *token.FileSet
	stateId   int
	labelId   int
	variables map[string]string
}

func (flt *Flattener) render(node ast.Node) string {
	var out bytes.Buffer
	format.Node(&out, flt.fset, node)
	return out.String()
}

func (flt *Flattener) addVariable(name, typ string) {
	if flt.variables == nil {
		flt.variables = make(map[string]string)
	}
	flt.variables[name] = typ
}

func (flt *Flattener) getLabelId() int {
	flt.labelId++
	return flt.labelId
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
		stateId := flt.getStateId()
		results := node.(*ast.ReturnStmt).Results
		if results == nil || len(results) != 1 {
			panic("Only supports a single result!")
		}
		returnValue := flt.render(results[0])
		return renderReturn(stateId, returnValue)

	case *ast.BlockStmt:
		builder := strings.Builder{}
		for _, stmt := range node.(*ast.BlockStmt).List {
			builder.WriteString(flt.flatten(stmt))
		}
		return builder.String()

	case *ast.IfStmt:
		ifStmt := node.(*ast.IfStmt)
		cond := flt.render(ifStmt.Cond)
		thenLabel := fmt.Sprintf("__then_%d", flt.getLabelId())
		elseLabel := fmt.Sprintf("__else_%d", flt.getLabelId())
		postLabel := fmt.Sprintf("__post_%d", flt.getLabelId())
		thenBody := flt.flatten(ifStmt.Body)
		elseBody := flt.flatten(ifStmt.Else)
		return renderIf(cond, thenLabel, thenBody, elseLabel, elseBody, postLabel)

	case *ast.ForStmt:
		forStmt := node.(*ast.ForStmt)
		if IsAllNil(forStmt.Init, forStmt.Cond, forStmt.Post) {
			return renderForever(fmt.Sprintf("__for_%d", flt.getLabelId()), flt.flatten(forStmt.Body))
		}
	case *ast.DeclStmt:
		name, typ := parseDeclStmt(node.(*ast.DeclStmt))
		if len(name) == 0 {
			break
		}
		flt.addVariable(name, typ)
		return fmt.Sprintf("// Original declaration of %s as %s\n", name, typ)
	case *ast.AssignStmt:
		// TODO: Actually support it! We currently can't declare vars in assignments.
		return fmt.Sprintln(strings.Replace(flt.render(node), ":=", "=", -1))
	case *ast.IncDecStmt:
		return fmt.Sprintln(flt.render(node))
	}
	return flt.showUnsupported(node)
}

func (flt *Flattener) showUnsupported(node ast.Node) string {
	var out bytes.Buffer
	ast.Fprint(&out, flt.fset, node, nil)
	if node != nil {
		return fmt.Sprintf("/* UNSUPPORTED: %v\n%v*/", flt.render(node), out.String())
	}
	return fmt.Sprintf("// UNSUPPORTED: nil\n")
}

func parseDeclStmt(decl *ast.DeclStmt) (name, typ string) {
	genDecl, ok := decl.Decl.(*ast.GenDecl)
	if !ok {
		return "", ""
	}
	name = genDecl.Specs[0].(*ast.ValueSpec).Names[0].Name
	typ = genDecl.Specs[0].(*ast.ValueSpec).Type.(*ast.Ident).Name
	return
}

//func (flt *Flattener) parseAssignStmt(assign *ast.AssignStmt, addVar func(name, typ string)) string {
//	for i, expr := range assign.Lhs {
//		ident, ok := expr.(*ast.Ident)
//		if !ok {
//			return flt.showUnsupported(assign)
//		}
//		if ident.Obj.Decl != nil {
//			return flt.showUnsupported(assign)
//		}
//		addVar(ident.Name, )
//	}
//}

func IsAllNil(things ...any) bool {
	for _, thing := range things {
		if thing != nil {
			return false
		}
	}
	return true
}

func (flt *Flattener) FlattenFunction(fd *ast.FuncDecl) string {

	_, name, _ := strings.Cut(fd.Name.Name, "_")
	body := flt.flatten(fd.Body)

	// This is a terrible hack!
	signature := flt.render(fd.Type)
	_, after, _ := strings.Cut(signature, "(")
	params, _, _ := strings.Cut(after, ")")

	return renderFunction(name, params, body, flt.stateId, flt.variables)
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

type Set[T comparable] struct {
	set map[T]struct{}
}

func (set *Set[T]) Add(value T) {
	set.set[value] = struct{}{}
}

func (set *Set[T]) Contains(value T) bool {
	_, exists := set.set[value]
	return exists
}

func (set *Set[T]) ExcludeSlice(slice []T) {
	for _, value := range slice {
		delete(set.set, value)
	}
}

func (set Set[T]) Iter() itertools.Iterator[T] {
	return itertools.FromSet(set.set)
}

func First[T any](first T, tail ...any) T {
	return first
}

func Second[T any](first any, second T, tail ...any) T {
	return second
}

func SetDifference[T comparable](a []T, b []T) (diff []T) {
	setB := make(map[T]struct{}, len(b))
	for _, item := range b {
		setB[item] = struct{}{}
	}
	for _, item := range a {
		if _, ok := setB[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func getFunctionDeclarations(typesInfo *types.Info, fdecl *ast.FuncDecl, fset *token.FileSet) {
	// Get function scope
	scope, exists := typesInfo.Scopes[fdecl.Type]
	if !exists {
		log.Fatal("No scope for function!")
	}

	// Get all parameter names
	var paramNames []string
	for _, field := range fdecl.Type.Params.List {
		for _, ident := range field.Names {
			paramNames = append(paramNames, ident.Name)
		}
	}

	// Get all result names
	var resultNames []string
	for _, field := range fdecl.Type.Results.List {
		for _, ident := range field.Names {
			resultNames = append(resultNames, ident.Name)
		}
	}

	// Get the names of all variables declared in the function (not including nested scopes...)
	scopeVarNames := SetDifference(SetDifference(scope.Names(), paramNames), resultNames)

	//vid := 0
	//
	//for childIndex := 0; childIndex < scope.NumChildren() ; childIndex++ {
	//	child := scope.Child(childIndex)
	//	for _, name := range child.Names() {
	//		child.Lookup(name).
	//	}
	//}
	//if scope.NumChildren() > 0 {
	//	itertools.ForEach(itertools.FromSlice(funcDeclScope.Child(0).Names()), func(x string) { fmt.Println("___ ", x) })
	//}

	fmt.Println("Scoped Vars:", scopeVarNames)

	for _, name := range scopeVarNames {
		obj := scope.Lookup(name)
		fmt.Println("var", name, obj.Type().String())
	}
}

func main() {
	pkgpath := "C:\\Code\\Personal\\go-explore\\flattener"
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedFiles | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedSyntax, Tests: false}
	pkgs, err := packages.Load(cfg, pkgpath)
	if err != nil {
		log.Fatal(err) // load error
	}

	var pkg *packages.Package
	for _, curPkg := range pkgs {
		if curPkg.Types.Name() == "main" {
			pkg = curPkg
			break
		}
	}
	if pkg == nil {
		log.Fatal("Couldn't find main package!")
	}

	var af *ast.File
	for _, f := range pkg.Syntax {
		//fmt.Println(f.Name.Name)
		if f.Name.Name == "main" {
			af = f
			break
		}
	}
	if af == nil {
		log.Fatal("Couldn't find file")
	}

	fset := pkg.Fset

	funcDecls := collectFuncDecls(af, false)

	//for key, val := range pkg.TypesInfo.Scopes {
	//	switch key.(type) {
	//	case *ast.FuncDecl:
	//		fmt.Println(key.(*ast.FuncDecl).Name.Name, val)
	//		continue
	//	}
	//	//fmt.Println(key, val)
	//}

	for _, fdecl := range funcDecls {
		funcDeclScope := pkg.TypesInfo.Scopes[fdecl.Type]
		if funcDeclScope == nil {
			log.Fatal("Shit.")
		}

		if !strings.HasPrefix(fdecl.Name.Name, "generate_") {
			continue
		}

		fmt.Println("X", fdecl.Name.Name)
		for i, param := range fdecl.Type.Params.List {
			fmt.Println(" :  ", i, param.Names)
		}
		for _, name := range funcDeclScope.Names() {
			fmt.Println("    ", name, funcDeclScope.Lookup(name))
		}

		if funcDeclScope.NumChildren() > 0 {
			itertools.ForEach(itertools.FromSlice(funcDeclScope.Child(0).Names()), func(x string) { fmt.Println("___ ", x) })
		}

		getFunctionDeclarations(pkg.TypesInfo, fdecl, pkg.Fset)
		//fmt.Println(fdecl.Name, funcDeclScope.Names())
	}

	var out bytes.Buffer

	out.WriteString("package main\n\n// AUTOGENERATED!\n\n")

	itertools.ForEach(itertools.Map(
		func(fdecl *ast.FuncDecl) string {
			flt := Flattener{fset: fset}
			return flt.FlattenFunction(fdecl)
		},
		itertools.FilterIn(
			itertools.FromSlice(funcDecls),
			func(fdecl *ast.FuncDecl) bool {
				return strings.HasPrefix(fdecl.Name.Name, "generate_")
			},
		),
	),
		func(f string) {
			out.WriteByte('\n')
			out.WriteString(f)
		},
	)

	//formattedResult, err := format.Source(out.Bytes())
	//if err != nil {
	//	fmt.Println(out.String())
	//	panic(err)
	//}

	formattedResult := out.Bytes()

	const target = "generators_gen.go"
	if err := os.WriteFile(target, formattedResult, 0644); err != nil {
		log.Fatal(err)
	}

}
