package analyzer

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestRangeValRef(t *testing.T) {
	wd, err := os.Getwd()
	require.Nil(t, err)
	testdata := filepath.Join(wd, "testdata")
	analysistest.Run(t, testdata, getRangeVarRefAnalyzer(), "range_val_ref")
}
