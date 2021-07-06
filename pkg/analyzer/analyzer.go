package analyzer

import (
	"sync"

	"golang.org/x/tools/go/analysis"
)

type analyzers struct {
	sync.Mutex
	analyzers []*analysis.Analyzer
}

var a analyzers

func registerAnalyzer(n *analysis.Analyzer) {
	a.Lock()
	defer a.Unlock()
	a.analyzers = append(a.analyzers, n)
}

func GetAnalyzers() []*analysis.Analyzer {
	a.Lock()
	defer a.Unlock()
	return a.analyzers
}
