// Package analyzer contains the code which carries out our linting
// check.
package analyzer

import (
	"go/ast"

	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"golang.org/x/tools/go/analysis"
)

// Analyzer configures our function.
var Analyzer = &analysis.Analyzer{
	Name:     "alphavet",
	Doc:      "Checks that functions are ordered alphabetically within packages.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// run is driven by the analysis framework, via the configuration above.
func run(pass *analysis.Pass) (interface{}, error) {

	// We're going to analyze things
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// But only functions
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// The previous function name, on a per-file basis
	prev := make(map[string]string)

	inspector.Preorder(nodeFilter, func(node ast.Node) {

		// Get the node, in the right type
		funcDecl := node.(*ast.FuncDecl)

		// Get the name of the function
		name := funcDecl.Name.Name

		// Get the name of the file we're analyzing
		//
		// Since we analyze based on packages, but packages
		// may contain multiple source files
		file := pass.Fset.File(node.Pos()).Name()

		// ignore special cases
		if ( name == "" || name == "init" || name == "main" ) {
			return
		}

		// Is this function out of order?
		if name != "" && ( name < prev[file] ) {

			pass.Reportf(node.Pos(),
				"function %s should have been before %s",
				name, prev[file])
		}

		prev[file] = name
	})

	return nil, nil
}
