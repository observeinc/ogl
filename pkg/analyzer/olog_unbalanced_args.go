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

func init() {
	registerAnalyzer(getOlogUnbalancedArgsAnalyzer())
}

// getOlogUnbalancedArgsAnalyzer builds and returns an Analyzer for the olog-unbalanced-args check.
func getOlogUnbalancedArgsAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "ologUnbalancedArgs",
		Doc:      "Checks that olog arguments are matched string/value pairs",
		Run:      checkOlogUnbalancedArgs,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

func isOlogCall(s *ast.CallExpr) bool {
	if fn, is := s.Fun.(*ast.SelectorExpr); is {
		if ident, is := fn.X.(*ast.Ident); is && ident.Name == "olog" {
			// olog.Info/Error
			if fn.Sel.Name == "Info" || fn.Sel.Name == "Error" {
				// olog.Info/Error
				return true
			}
		} else if call, is := fn.X.(*ast.CallExpr); is {
			// olog.V(n).Info/Error
			if sel, is := call.Fun.(*ast.SelectorExpr); is {
				if x, is := sel.X.(*ast.Ident); is && x.Name == "olog" && sel.Sel.Name == "V" {
					if fn.Sel.Name == "Info" || fn.Sel.Name == "Error" {
						// olog.Info/Error
						return true
					}
				}
			}
		}
	}
	return false
}

func typeIsStringOrInterfaceArray(argt types.TypeAndValue) bool {
	return false // TODO
}

func ologCallHasInvalidCount(s *ast.CallExpr, ti *types.Info) bool {
	// first: It may be an unpack -- if so, allow it
	if len(s.Args) == 2 {
		if s.Ellipsis == token.NoPos {
			// not an ellipsis unpack
			return true
		}
		return false
	}
	if (len(s.Args) & 1) != 1 {
		// not a string, key, value, key, value, ... matched set
		return true
	}
	return false
}

func ologCallArgumentIsNotString(s *ast.CallExpr, ti *types.Info) bool {
	if len(s.Args) < 3 {
		// can't check when there's not args
		return false
	}
	for i := 1; i < len(s.Args); i += 2 {
		if typ, has := ti.Types[s.Args[i]]; has {
			if bt, is := typ.Type.(*types.Basic); is {
				if bt.Kind() != types.String {
					// this is not a string! (retyped strings are OK)
					return true
				}
			} else {
				// this is not a string, because it's not a basic type
				return true
			}
		}
	}
	// passed all checks, so no odd argument is not a string
	return false
}

// checkOlogUnbalancedArgs checks if there's a call to an olog function that
// requires string/value pairs, that doesn't get that.
func checkOlogUnbalancedArgs(pass *analysis.Pass) (interface{}, error) {
	visitor := func(node ast.Node) {
		switch s := node.(type) {
		case *ast.CallExpr:
			if isOlogCall(s) {
				if ologCallHasInvalidCount(s, pass.TypesInfo) {
					pass.Reportf(s.Pos(), "olog key is missing value")
				} else if ologCallArgumentIsNotString(s, pass.TypesInfo) {
					pass.Reportf(s.Pos(), "olog key is not a string")
				}
			}
		default:
			// this shouldn't happen, the filter below shouldn't let this happen
			panic(fmt.Sprintf("Unexpected statement of type %T received", s))
		}
	}

	nodeFilter := []ast.Node{
		(*ast.CallExpr)(nil),
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder(nodeFilter, visitor)
	return nil, nil
}
