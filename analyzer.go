// Package analyzer contains the code which carries out our linting
// check.
package analyzer

import (
	"go/ast"
	"go/token"
	"sort"
	"strings"

	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"

	"golang.org/x/tools/go/analysis"
)

// exclude is the list of files to exclude
var exclude stringSetFlag

func init() {
	exclude.Set("init,main")
	Analyzer.Flags.Var(&exclude, "exclude",
		"comma-separated list of functions to exclude from the ordering test.")

}

// Analyzer configures our function.
var Analyzer = &analysis.Analyzer{
	Name:     "alphavet",
	Doc:      "Checks that functions are ordered alphabetically within packages.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
}

// Function represents a function definition which
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

// excludeFunction determines if the given function should be excluded
// from our alphabetical constraints.
//
// Typically we'd exclude "init", and "main".  Optionally the user might
// disable this, or add more exclusions.
func excludeFunction(name string) bool {

	_, found := exclude[name]
	return found
}

// findFunctions will use the analysis package to return a list of
// all the functions which were found in the specified package(s).
func findFunctions(pass *analysis.Pass) []*Function {

	// We're going to analyze things
	inspector := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	// But only functions
	nodeFilter := []ast.Node{
		(*ast.FuncDecl)(nil),
	}

	// The a list of all the functions we've found will be returned
	// to the caller.
	seen := []*Function{}

	// Add to the list all the function definitions we encounter, in order.
	inspector.Preorder(nodeFilter, func(node ast.Node) {

		// Get the node, in the right type
		funcDecl := node.(*ast.FuncDecl)

		if excludeFunction(funcDecl.Name.Name) {
			return
		}

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

	// Return the list
	return seen
}

// run is driven by the analysis framework and gets passed instances
// of function definitions.
func run(pass *analysis.Pass) (interface{}, error) {

	//
	// Find the functions in the package(s) we've been
	// asked to lint/vet.
	//
	seen := findFunctions(pass)

	//
	// We want to find the unique filename we've seen.
	//
	// Remember we're invoked on "packages", but packages
	// may be implemented by multiple files.
	//
	files := make(map[string]bool)
	for _, fnc := range seen {
		files[fnc.File] = true
	}

	sortedFiles := make([]string, 0, len(files))

	for k := range files {
		sortedFiles = append(sortedFiles, k)
	}
	sort.Strings(sortedFiles)

	//
	// Now we have a list of unique filenames.
	//
	// Process each one.
	//
	for _, fn := range sortedFiles {

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

type stringSetFlag map[string]bool

func (ss *stringSetFlag) String() string {
	var items []string
	for item := range *ss {
		items = append(items, item)
	}
	sort.Strings(items)
	return strings.Join(items, ",")
}

func (ss *stringSetFlag) Set(s string) error {
	m := make(map[string]bool) // clobber previous value
	if s != "" {
		for _, name := range strings.Split(s, ",") {
			if name == "" {
				continue // TODO: report error? proceed?
			}
			m[name] = true
		}
	}
	*ss = m
	return nil
}
