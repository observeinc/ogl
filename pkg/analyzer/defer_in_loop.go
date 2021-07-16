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
	// max number of nested loops + function literals we expect. This is just for preallocating memory.
	expectedMaxScopes = 5
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
	type scopeInfo struct {
		isLoop     bool
		start, end token.Pos
	}
	scopes := make([]scopeInfo, expectedMaxScopes)
	isInLoop := func(pos token.Pos) bool {
		for i := len(scopes) - 1; i >= 0; i-- {
			if pos > scopes[i].start && pos < scopes[i].end {
				return scopes[i].isLoop
			}
			scopes = scopes[:i]
		}
		return false
	}
	addFuncScope := func(start, end token.Pos) {
		scopes = append(scopes, scopeInfo{false, start, end})
	}

	addLoopScope := func(start, end token.Pos) {
		scopes = append(scopes, scopeInfo{true, start, end})
	}
	visitor := func(node ast.Node) {
		switch s := node.(type) {
		case *ast.FuncLit:
			addFuncScope(s.Pos(), s.End())
		case *ast.ForStmt:
			addLoopScope(s.Pos(), s.End())
		case *ast.RangeStmt:
			addLoopScope(s.Pos(), s.End())
		case *ast.DeferStmt:
			if isInLoop(s.Pos()) {
				pass.Reportf(s.Pos(), "found defer statement in loop")
			}
		default:
			// this shouldn't happen, the filter below shouldn't let this happen
			panic(fmt.Sprintf("Unexpected statement of type %T received", s))
		}
	}

	nodeFilter := []ast.Node{
		(*ast.FuncLit)(nil),
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
		(*ast.DeferStmt)(nil),
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder(nodeFilter, visitor)
	return nil, nil
}
