package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOlogUnbalancedArgs(t *testing.T) {
	wd, err := os.Getwd()
	require.Nil(t, err)
	testdata := filepath.Join(wd, "testdata")
	analysistest.Run(t, testdata, getOlogUnbalancedArgsAnalyzer(), "olog_unbalanced_args")
}
