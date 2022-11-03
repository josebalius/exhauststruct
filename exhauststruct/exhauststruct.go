package exhauststruct

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:      "exhauststruct",
	Doc:       "it checks that all struct fields are initialized",
	Run:       run,
	Requires:  []*analysis.Analyzer{inspect.Analyzer},
	FactTypes: []analysis.Fact{&structFact{}},
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	structs := findStructs(inspect, pass.TypesInfo)
	exportFacts(pass, structs)

	return structCheck(pass, inspect)
}

func exportFacts(pass *analysis.Pass, structs []structType) {
	for _, s := range structs {
		pass.ExportObjectFact(s.factObject(), &structFact{fields: s.fields})
	}
}

func structCheck(pass *analysis.Pass, inspect *inspector.Inspector) (interface{}, error) {
	structChecker := makeStructChecker(pass)
	inspect.Nodes([]ast.Node{&ast.CompositeLit{}}, func(n ast.Node, push bool) bool {
		return structChecker(n, push)
	})
	return nil, nil
}

func makeStructChecker(pass *analysis.Pass) func(ast.Node, bool) bool {
	return func(n ast.Node, push bool) bool {
		l := n.(*ast.CompositeLit)

		st := structTypeFromCompositeLit(l, pass.TypesInfo)
		fact, ok := importFact(pass, st)
		if !ok {
			return true
		}

		usedFields := compositeLitFields(l)
		unusedFields := make([]string, 0, len(fact.fields))
		for _, f := range fact.fields {
			if _, ok := usedFields[f]; !ok {
				unusedFields = append(unusedFields, f)
			}
		}

		if len(unusedFields) > 0 {
			pass.Reportf(l.Pos(), "uninitialized struct fields: %v", strings.Join(unusedFields, ", "))
		}

		return true
	}
}

func compositeLitFields(l *ast.CompositeLit) map[string]struct{} {
	fields := make(map[string]struct{})
	for _, elt := range l.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		fields[fmt.Sprintf("%s", kv.Key)] = struct{}{}
	}
	return fields
}

func importFact(pass *analysis.Pass, possibleStructType structType) (structFact, bool) {
	var f structFact
	found := pass.ImportObjectFact(possibleStructType.factObject(), &f)
	return f, found
}

func findStructs(inspect *inspector.Inspector, info *types.Info) []structType {
	var result []structType

	inspect.Preorder([]ast.Node{&ast.GenDecl{}}, func(n ast.Node) {
		gendecl := n.(*ast.GenDecl)
		if gendecl.Doc == nil || len(gendecl.Specs) != 1 {
			return
		}
		if exhaust := exhaustStruct(gendecl); !exhaust {
			return
		}

		spec := gendecl.Specs[0].(*ast.TypeSpec)
		specStructType, ok := spec.Type.(*ast.StructType)
		if !ok || specStructType.Fields == nil {
			return
		}

		var fields []string
		for _, f := range specStructType.Fields.List {
			fields = append(fields, fmt.Sprintf("%s", f.Names[0]))
		}

		st := structTypeFromGenDecl(gendecl, info)
		st.fields = fields
		result = append(result, st)
	})

	return result
}

func exhaustStruct(gendecl *ast.GenDecl) bool {
	for _, c := range gendecl.Doc.List {
		if c.Text == "//lint:exhauststruct" {
			return true
		}
	}
	return false
}

type structType struct {
	*types.TypeName

	fields []string
}

func structTypeFromCompositeLit(l *ast.CompositeLit, info *types.Info) structType {
	t := info.Types[l.Type]
	tagType, ok := t.Type.(*types.Named)
	if !ok {
		return structType{}
	}

	return structType{TypeName: tagType.Obj()}
}

func structTypeFromGenDecl(g *ast.GenDecl, info *types.Info) structType {
	spec := g.Specs[0].(*ast.TypeSpec)
	obj := info.Defs[spec.Name]
	named := obj.Type().(*types.Named)
	typeName := named.Obj()

	return structType{TypeName: typeName}
}

func (et structType) String() string           { return et.TypeName.String() }
func (et structType) scope() *types.Scope      { return et.TypeName.Parent() }
func (et structType) factObject() types.Object { return et.TypeName }

type structFact struct {
	fields []string
}

func (s *structFact) AFact()         {}
func (s *structFact) String() string { return "structFact" }
