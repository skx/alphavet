// Package analyzer contains the code which carries out our linting
// check.
package analyzer

import (
	"go/ast"
	"go/token"

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

// Type Function represents a function definition which
// has been encountered when scanning the source code.
//
// We want to differentiate between function calls that have
// receivers, and those that don't so we store some extra
// details here.
type Function struct {
	// The file containing the function.
	File string

	// The name of the function.
	Name string

	// The position within our source where this object was.
	Position token.Pos

	// The receiver, if any, for the function.
	Receiver string
}

// run is driven by the analysis framework and gets passed instances
// of function definitions.
func run(pass *analysis.Pass) (interface{}, error) {

	// We're going to analyze things
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// But only functions
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// Create a list of all the functions we've found
	seen := []*Function{}

	// Now update that list with all the function definitions
	// we encounter, in order.
	inspector.Preorder(nodeFilter, func(node ast.Node) {

		// Get the node, in the right type
		funcDecl := node.(*ast.FuncDecl)

		// Build up any (optional) receiver
		recv := ""

		// Is there a receiver?
		if funcDecl.Recv != nil {

			for _, x := range funcDecl.Recv.List {

				if x.Names != nil {
					for _, y := range x.Names {
						recv += y.Name
					}
				}
			}
		}

		// The entry for the function, with appropriate data
		tmp := &Function{
			File:     pass.Fset.File(node.Pos()).Name(),
			Name:     funcDecl.Name.Name,
			Position: node.Pos(),
			Receiver: recv,
		}

		// Save it away
		seen = append(seen, tmp)
	})

	//
	// Now we've processed all the functions.
	//
	// We want to build up a list of files we've seen.
	//
	// Remember we're invoked on "packages", but packages
	// may contain multiple files.
	//
	files := make(map[string]bool)
	for _, fnc := range seen {
		files[fnc.File] = true
	}

	//
	// Now we have a list of unique filenames.
	//
	// Process each one.
	//
	for fn := range files {

		//
		// Function names we've seen that have a receiver
		//
		obj := make(map[string][]*Function)

		//
		// Function names with no receiver
		//
		raw := []*Function{}

		//
		// Append the function for this file
		// to the appropriate slice
		//
		for _, ent := range seen {

			// Not for this file?  Skip for now
			if ent.File != fn {
				continue
			}

			//
			// Record the method, with either a receiver
			// or not.
			//
			if ent.Receiver == "" {
				raw = append(raw, ent)
			} else {
				obj[ent.Receiver] = append(obj[ent.Receiver], ent)
			}
		}

		//
		// Now we have a list of methods with receivers
		// and not
		//

		//
		// Test each receiver
		//
		seen := make(map[string]*Function)
		for _, entries := range obj {

			for _, ent := range entries {

				if (seen[ent.Receiver] != nil) && (ent.Name < seen[ent.Receiver].Name) {
					pass.Reportf(ent.Position,
						"function %s on receiver %s should have been before %s",
						ent.Name, ent.Receiver, seen[ent.Receiver].Name)
				}
				seen[ent.Receiver] = ent
			}
		}

		//
		// Then test those functions with no receiver
		//
		prev := ""

		for _, r := range raw {
			if (prev != "") && (r.Name < prev) {
				pass.Reportf(r.Position,
					"function %s should have been before %s",
					r.Name, prev)
			}
			prev = r.Name
		}
	}

	return nil, nil
}
