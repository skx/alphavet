// Package main contains the main driver, which launches our linter.
package main

import (
	"flag"

	"github.com/skx/alphavet/analyzer"
	"golang.org/x/tools/go/analysis/singlechecker"
)

// entry-point
func main() {

	// We don't use any flags, but "go vet" might have this passed.
	flag.Bool("unsafeptr", false, "")

	singlechecker.Main(analyzer.Analyzer)
}
