package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

const (
	// max number of range variables we expect. This is just for preallocating memory.
	// 2 range vars per loop * 5 levels of nesting should be good enough for real world code.
	expectedMaxRangeVars = 10
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
		name       string
		scopeStart token.Pos
		scopeEnd   token.Pos
	}
	rangeVars := make([]rangeVar, 0, expectedMaxRangeVars)

	check := func(a ast.Expr, pos token.Pos, reporter func(id *ast.Ident)) {
		for {
			switch id := a.(type) {
			case *ast.Ident:
				for i := len(rangeVars) - 1; i >= 0; i-- {
					r := rangeVars[i]
					if id.Name == r.name && pos > r.scopeStart && pos < r.scopeEnd {
						reporter(id)
						break
					}
				}
				return
			case *ast.SelectorExpr:
				if t, ok := pass.TypesInfo.Types[id.X]; ok {
					if _, isPointer := t.Type.(*types.Pointer); isPointer {
						return
					}
				}
				a = id.X
			case *ast.ParenExpr:
				a = id.X
			case *ast.IndexExpr:
				// skip, can assign to or take reference of items in nested
				// containers via range variables.
				return
			case *ast.StarExpr:
				// skip: a valid star expression obviously has a pointer.
				return
			default:
				panic(fmt.Sprintf("Unknown expr type: %T", id))
			}
		}
	}
	visitor := func(node ast.Node) {
		// cleanup rangeVars that node is out of bounds of
		for i := 0; i < len(rangeVars); i++ {
			p := node.Pos()
			if p < rangeVars[i].scopeStart || p > rangeVars[i].scopeEnd {
				if i == len(rangeVars)-1 {
					rangeVars = rangeVars[:i]
				} else {
					n := len(rangeVars) - 1
					rangeVars[i] = rangeVars[n]
					rangeVars = rangeVars[:n]
					i--
				}
			}
		}
		switch s := node.(type) {
		case *ast.RangeStmt:
			if key, ok := s.Key.(*ast.Ident); ok && key.Name != "_" {
				rangeVars = append(rangeVars, rangeVar{key.Name, s.Pos(), s.End()})
			}
			if val, ok := s.Value.(*ast.Ident); ok && val.Name != "_" {
				rangeVars = append(rangeVars, rangeVar{val.Name, s.Pos(), s.End()})
			}
		case *ast.UnaryExpr:
			if s.Op != token.AND {
				return
			}
			if _, isComposite := s.X.(*ast.CompositeLit); isComposite {
				// address of a literal, e.g. &myStruct{}
				return
			}
			check(s.X, s.Pos(), func(id *ast.Ident) {
				pass.Reportf(s.Pos(),
					"taking reference of range variable %s or its field",
					id.Name)
			})
		case *ast.IncDecStmt:
			check(s.X, s.Pos(), func(id *ast.Ident) {
				pass.Reportf(s.Pos(),
					"modifying range variable %s or its field",
					id.Name)
			})
		case *ast.AssignStmt:
			for i := 0; i < len(s.Lhs); i++ {
				check(s.Lhs[i], s.Pos(), func(id *ast.Ident) {
					pass.Reportf(s.Pos(),
						"modifying range variable %s or its field",
						id.Name)
				})
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
		(*ast.IncDecStmt)(nil),
		(*ast.AssignStmt)(nil),
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder(nodeFilter, visitor)
	return nil, nil
}
