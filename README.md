# alphavet

This is a simple linter which is designed to report upon functions which are not implemented in alphabetical order within files.

The motivation behind this tool was twofold:

* I find it easier to navigate functions if they are ordered alphabetically.
* Once I realized a linter, driven by "`go vet`", could be named "alphavet" I couldn't resist the temptation to hack it up.
  * Even though this could just has easily have been a portable Perl script.


## Installation

If you have a working golang toolset you should be able to install by:

```sh
go install github.com/skx/alphavet/cmd/alphavet@latest
```



## Usage

The linter is designed to be driven by `go vet` like so:

```sh
$ go vet -vettool=$(which alphavet) ./...
```

## Sample Output

Sample output would look something like this:

```
$ go vet -vettool=$(which alphavet) ./...
# github.com/skx/gobasic/builtin
./builtin.go:67:1: function Get should have been before Register
./misc_test.go:21:1: function LineEnding should have been before StdInput
./misc_test.go:29:1: function StdError should have been before StdOutput
./misc_test.go:33:1: function Data should have been before StdError

```
