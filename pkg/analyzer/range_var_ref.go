package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	expectedMaxNesting = 5
)

func init() {
	registerAnalyzer(getRangeVarRefAnalyzer())
}

// getRangeVarRefAnalyzer builds and returns an Analyzer for the defer-in-loop check.
func getRangeVarRefAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "rangeVarRef",
		Doc:      "Checks that references to range variables aren't taken",
		Run:      checkRangeVarRef,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// checkRangeVarRef checks if a defer statement exists anywhere in a loop and flags it as an issue.
func checkRangeVarRef(pass *analysis.Pass) (interface{}, error) {
	type rangeVar struct {
		name     string
		scopeEnd token.Pos
	}
	rangeVars := make([]rangeVar, 0, expectedMaxNesting)
	visitor := func(node ast.Node) {
		// cleanup rangeVars that node is out of bounds of
		for i, r := range rangeVars {
			if r.scopeEnd < node.Pos() {
				rangeVars = rangeVars[:i]
				break
			}
		}
		switch s := node.(type) {
		case *ast.RangeStmt:
			if key, ok := s.Key.(*ast.Ident); ok && key.Name != "_" {
				rangeVars = append(rangeVars, rangeVar{key.Name, s.End()})
			}
			if val, ok := s.Value.(*ast.Ident); ok && val.Name != "_" {
				rangeVars = append(rangeVars, rangeVar{val.Name, s.End()})
			}
		case *ast.UnaryExpr:
			if s.Op != token.AND {
				return
			}
			var id *ast.Ident
			switch i := s.X.(type) {
			case *ast.Ident:
				id = i
			case *ast.SelectorExpr:
				id, _ = i.X.(*ast.Ident)
			}
			if id == nil {
				return
			}

			for _, r := range rangeVars {
				if id.Name == r.name && s.Pos() < r.scopeEnd {
					pass.Reportf(s.Pos(),
						"taking reference of range variable %s or its field",
						id.Name)
					break
				}
			}
		default:
			// this shouldn't happen, the filter below shouldn't let this happen
			panic(fmt.Sprintf("Unexpected statement of type %T received", s))
		}
		return
	}

	nodeFilter := []ast.Node{
		(*ast.RangeStmt)(nil),
		(*ast.UnaryExpr)(nil),
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder(nodeFilter, visitor)
	return nil, nil
}
