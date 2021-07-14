package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	expectedMaxLoopNesting = 5
)

func init() {
	registerAnalyzer(getRangeVarRefAnalyzer())
}

// getRangeVarRefAnalyzer builds and returns an Analyzer for the range-variable-reference check.
func getRangeVarRefAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "rangeVarRef",
		Doc:      "Checks that references to range variables aren't taken",
		Run:      checkRangeVarRef,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// checkRangeVarRef checks if a range variable's reference is taken.
func checkRangeVarRef(pass *analysis.Pass) (interface{}, error) {
	type rangeVar struct {
		name     string
		scopeEnd token.Pos
	}
	rangeVars := make([]rangeVar, 0, expectedMaxLoopNesting)

	checkIfRangeVar := func(id *ast.Ident, pos token.Pos) {
		for i := len(rangeVars) - 1; i >= 0; i-- {
			r := rangeVars[i]
			if id.Name == r.name && pos < r.scopeEnd {
				pass.Reportf(pos,
					"taking reference of range variable %s or its field",
					id.Name)
				break
			}
		}
	}

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

			a := s.X
			for {
				switch id := a.(type) {
				case *ast.Ident:
					checkIfRangeVar(id, s.Pos())
					return
				case *ast.SelectorExpr:
					if t, ok := pass.TypesInfo.Types[id.X]; ok {
						if _, isPointer := t.Type.(*types.Pointer); isPointer {
							return
						}
					}
					a = id.X
				default:
					panic(fmt.Sprintf("Unknown expr type: %s", spew.Sdump(id)))
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
