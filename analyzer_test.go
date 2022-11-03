package analyzer_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skx/alphavet"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAll(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get wd: %s", err)
	}

	testdata := filepath.Join(wd, "testdata")
	analysistest.Run(t, testdata, analyzer.Analyzer, "p")
}
