package analyzer

import (
	"fmt"
	"go/ast"
	"go/token"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func init() {
	registerAnalyzer(getDeferInLoopAnalyzer())
}

// getDeferInLoopAnalyzer builds and returns an Analyzer for the defer-in-loop check.
func getDeferInLoopAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "deferInLoop",
		Doc:      "Checks that there are no defer statements in loops",
		Run:      checkDeferInLoop,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// checkDeferInLoop checks if a defer statement exists anywhere in a loop and flags it as an issue.
func checkDeferInLoop(pass *analysis.Pass) (interface{}, error) {
	loopEnd := token.Pos(0)
	visitor := func(node ast.Node) {
		switch s := node.(type) {
		case *ast.ForStmt:
			if s.End() > loopEnd {
				loopEnd = s.End()
			}
		case *ast.RangeStmt:
			if s.End() > loopEnd {
				loopEnd = s.End()
			}
		case *ast.DeferStmt:
			if s.Pos() < loopEnd {
				pass.Reportf(s.Pos(), "found defer statement in loop")
			}
		default:
			// this shouldn't happen, the filter below shouldn't let this happen
			panic(fmt.Sprintf("Unexpected statement of type %T received", s))
		}
	}

	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
		(*ast.DeferStmt)(nil),
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder(nodeFilter, visitor)
	return nil, nil
}
