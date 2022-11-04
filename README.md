[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/skx/alphavet)
[![Go Report Card](https://goreportcard.com/badge/github.com/skx/alphavet)](https://goreportcard.com/report/github.com/skx/alphavet)
[![license](https://img.shields.io/github/license/skx/alphavet.svg)](https://github.com/skx/alphavet/blob/master/LICENSE)



* [alphavet](#alphavet)
  * [Installation](#installation)
  * [Usage](#usage)
* [Sample Output](#sample-output)
* [Github Setup](#github-setup)
* [Bug reports?](#bug-reports?)



# alphavet

This is a simple linter which is designed to report upon functions which are not implemented in alphabetical order within files.

The motivation behind this tool was twofold:

* I find it easier to navigate functions if they are ordered alphabetically.
  * Most IDEs offer a tree/outline view which is ordered alphabetically, and the contents and the tree should match!
* Once I realized a linter, driven by "`go vet`", could be named "alphavet" I couldn't resist the temptation to hack it up.
  * Even though this could just has easily have been a portable Perl script.

**NOTE** That the two functions `init` and `main` are excluded from the alphabetical ordering requirement.

* If there is interest/demand I could make a similiar, optional, exclusion for `New` and `NewXXX` functions in the future.



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


# Sample Output

Sample output would look something like this:

```
$ go vet -vettool=$(which alphavet) ./...
# github.com/skx/gobasic/builtin
./builtin.go:67:1: function Get should have been before Register
./misc_test.go:21:1: function LineEnding should have been before StdInput
./misc_test.go:29:1: function StdError should have been before StdOutput
./misc_test.go:33:1: function Data should have been before StdError

```



# Github Setup

This repository is configured to run tests upon every commit, and when pull-requests are created/updated.  The testing is carried out via [.github/run-tests.sh](.github/run-tests.sh) which is used by the [github-action-tester](https://github.com/skx/github-action-tester) action.



# Bug reports?

Please do feel free to report any issues you see with the code, or the results.

Feature requests are also welcome, although I'd prefer to avoid having excessive flags.




Steve
--
