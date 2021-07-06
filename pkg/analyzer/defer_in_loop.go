package analyzer

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

func init() {
	registerAnalyzer(getDeferInloopAnalyzer())
}

// getDeferInloopAnalyzer builds and returns an Analyzer for the defer-in-loop check.
func getDeferInloopAnalyzer() *analysis.Analyzer {
	return &analysis.Analyzer{
		Name:     "deferInLoop",
		Doc:      "Checks that there are no defer statements in loops",
		Run:      checkDeferInLoop,
		Requires: []*analysis.Analyzer{inspect.Analyzer},
	}
}

// checkDeferInLoop checks if a defer statement exists anywhere in a loop and flags it as an issue.
func checkDeferInLoop(pass *analysis.Pass) (interface{}, error) {
	visitor := func(node ast.Node) {
		var stmts []ast.Stmt
		switch s := node.(type) {
		case *ast.ForStmt:
			stmts = make([]ast.Stmt, len(s.Body.List))
			copy(stmts, s.Body.List)
		case *ast.RangeStmt:
			stmts = make([]ast.Stmt, len(s.Body.List))
			copy(stmts, s.Body.List)
		default:
			// this shouldn't happen, the filter below shouldn't let this happen
			panic(fmt.Sprintf("Unexpected statement of type %T received", s))
		}
		for len(stmts) > 0 {
			switch s := stmts[0].(type) {
			case *ast.DeferStmt:
				pass.Reportf(s.Pos(), "found defer statement in loop")
			case *ast.IfStmt:
				stmts = append(stmts, s.Body)
				stmts = append(stmts, s.Else)
			case *ast.BlockStmt:
				stmts = append(stmts, s.List...)
			case *ast.TypeSwitchStmt:
				stmts = append(stmts, s.Body)
			case *ast.SwitchStmt:
				stmts = append(stmts, s.Body)
			case *ast.CaseClause:
				stmts = append(stmts, s.Body...)
			case *ast.SelectStmt:
				stmts = append(stmts, s.Body)
			case *ast.CommClause:
				stmts = append(stmts, s.Body...)
			default:
				// some other type of statement that we don't care about.
				// Note, nested for/range loops will get their own call to inspect.
			}
			stmts = stmts[1:]
		}
	}

	nodeFilter := []ast.Node{
		(*ast.ForStmt)(nil),
		(*ast.RangeStmt)(nil),
	}
	insp := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
	insp.Preorder(nodeFilter, visitor)
	return nil, nil
}
