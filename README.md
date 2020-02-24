# gomodfuzz [![GoDoc](https://godoc.org/github.com/codeactual/gomodfuzz?status.svg)](https://pkg.go.dev/mod/github.com/codeactual/gomodfuzz) [![Go Report Card](https://goreportcard.com/badge/github.com/codeactual/gomodfuzz)](https://goreportcard.com/report/github.com/codeactual/gomodfuzz) [![Build Status](https://travis-ci.org/codeactual/gomodfuzz.png)](https://travis-ci.org/codeactual/gomodfuzz)

gomodfuzz is a program which assists testing of Go program compatibility with 1.11+ module support.

It runs the input program with permutations of `GO111MODULE`, `GOFLAGS`, `GOPATH`, and working directory traits (`GOMOD`, under `GOPATH`, etc.).

## Use Case

Testing programs which rely on parts of the Go toolchain such as [golang.org/x/tools/go/packages](https://pkg.go.dev/mod/golang.org/x/tools/go/packages) to load packages.

It was originally made to assert that [aws-mockery](https://github.com/codeactual/aws-mockery) could load packages from as many file location scenarios as possible.

## Permutation values

- working directory's is inside a module's file tree (`WD` in the output)
  - `true`
  - `false`
- working directory's relationship to `GOPATH` (`IN_MODULE` in the output)
  - inside `GOPATH`
  - outside `GOPATH`
- `GO111MODULE`
  - `auto`
  - `off`
  - `on`
- `GOFLAGS`
  - empty
  - `-mod=vendor`
- `GOPATH`
  - empty
  - a path which will contain the working directory if the "working directory's relationship to `GOPATH`" permutation value is "inside `GOPATH`"
  - a path which will never contain the working directory

# Usage

> To install: `go get -v github.com/codeactual/gomodfuzz/cmd/gomodfuzz`

## Examples

> Usage:

```bash
gomodfuzz --help
```

> Basic test:

```bash
gomodfuzz -- /path/to/subject --subject_flag0  --subject_flag1 subject_arg0 subject_arg1
```

> Run subject command with a timeout of 10 seconds:

```bash
gomodfuzz --timeout 10 -- /path/to/subject
```

> Display verbose results (passes, full errors, etc.)

```bash
gomodfuzz -v -- /path/to/subject
```

# Development

## License

[Mozilla Public License Version 2.0](https://www.mozilla.org/en-US/MPL/2.0/) ([About](https://www.mozilla.org/en-US/MPL/), [FAQ](https://www.mozilla.org/en-US/MPL/2.0/FAQ/))

## Contributing

- Please feel free to submit issues, PRs, questions, and feedback.
- Although this repository consists of snapshots extracted from a private monorepo using [transplant](https://github.com/codeactual/transplant), PRs are welcome. Standard GitHub workflows are still used.
