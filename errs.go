// Package errs is a convenience wrapper for chaining multiple
// error returning functions.
//  var e errs.Group
//
//  // add couple of functions
//  e.Add(func() error { ... })
//  e.Add(func() error { ... })
//  e.Add(func() error { ... })
//  e.Defer(func() error { ... }) // executes after other functons
//  e.Final(func() { ... })       // executes even if error is returned.
//
//  // execute them
//  if err := e.Exec(); err != nil {
//      // handle error
//  }
package errs

// Group is a group of functions that returns an error.
// Empty value of Group is usable.
type Group struct {
	funcs  []func() error
	defers []func() error
	final  []func()
}

// Add adds f to the group of functions.
// Functions are executed FIFO.
func (g *Group) Add(f func() error) {
	g.funcs = append(g.funcs, f)
}

// Defer adds f to the group of deferred functions.
// Similar to Add, Defer can be called multiple times
// to add more defer functions.
// Defer functions are executed LIFO.
func (g *Group) Defer(f func() error) {
	g.defers = append([]func() error{f}, g.defers...)
}

// Final is the function that is guaranteed to be executed
// even if an error is returned.
func (g *Group) Final(f func()) {
	g.final = append(g.final, f)
}

// Exec runs all functions then defer functions, stops on the
// first that errored and returns the error occurred.
// If no error is encountered, returns nil.
func (g Group) Exec() error {
	defer func() {
		for _, f := range g.final {
			f()
		}
	}()

	for _, f := range append(g.funcs, g.defers...) {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
