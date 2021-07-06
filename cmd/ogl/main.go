package main

import (
	"github.com/observeinc/ogl/pkg/analyzer"

	"golang.org/x/tools/go/analysis/multichecker"
)

func main() {
	multichecker.Main(analyzer.GetAnalyzers()...)
}
