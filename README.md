# gotchet [![build](https://github.com/alexbakker/gotchet/actions/workflows/build.yml/badge.svg)](https://github.com/alexbakker/gotchet/actions/workflows/build.yml)

__gotchet__ is a test report tool for Go. It can display test results in a TUI
and generate HTML reports.

## Usage

Use one of the following commands to get the right
[test2json](https://pkg.go.dev/cmd/test2json) output from your Go tests:

For the ``go test`` command:

```
go test -json -v=test2json ./...
```

For test binaries:

```
go tool test2json -t -p pkgname ./test-binary -test.v=test2json
```
