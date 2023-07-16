# gotchet

__gotchet__ is a test report tool for Go. It can display test reports in a TUI
and generate an HTML report.

## Usage

To get the right output format, use onf of the following commands:

For the ``go test`` command:

```
go test -json -v=test2json ./...
```

For test binaries: https://pkg.go.dev/cmd/test2json

```
go tool test2json -t -p pkgname ./test-binary -test.v=test2json
```
