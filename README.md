# errs
convenience wrapper for chaining multiple error returning functions.

[![Documentation](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/abiosoft/errs)

## Usage
```go
var e errs.Group

// add couple of functions
e.Add(func() error { ... })
e.Add(func() error { ... })
e.Add(func() error { ... })
e.Defer(func() error) // executes after other functons
e.Final(func())       // executes even if error is returned.

// execute them
if err := e.Exec(); err != nil {
    // handle error
}
```