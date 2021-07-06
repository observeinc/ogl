package main

import (
	"ogl/pkg/analyzer"

	"golang.org/x/tools/go/analysis"
)

// This must be defined and named as such for glangci-lint to be able to use this.
var AnalyzerPlugin analyzerPlugin

type analyzerPlugin struct{}

func (*analyzerPlugin) GetAnalyzers() []*analysis.Analyzer {
	return analyzer.GetAnalyzers()
}
