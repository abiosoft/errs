// Package errs is a convenience wrapper for chaining multiple
// error returning functions.
//  var e errs.Group
//
//  // add couple of functions
//  e.Add(func() error { ... })
//  e.Defer(func() { ... }) // executes after other functions
//  e.Add(func() error { ... })
//  e.Add(func() error { ... })
//  e.Final(func() { ... }) // executes even if error is returned
//
//  // execute them
//  if err := e.Exec(); err != nil {
//      // handle error
//  }
package errs

type fn struct {
	f func() error
	d bool // defer
}

// Group is a group of functions.
// Empty value of Group is usable.
type Group struct {
	funcs []fn
	final []func()
}

// Add adds f to the group of functions.
// Functions are executed FIFO.
func (g *Group) Add(f func() error) {
	g.funcs = append(g.funcs, fn{f: f})
}

// Defer adds f to the group of deferred functions.
// Similar to Add, Defer can be called multiple times
// to add more defer functions.
// Defer functions are executed LIFO.
func (g *Group) Defer(f func()) {
	g.funcs = append(g.funcs, fn{f: func() error { f(); return nil }, d: true})
}

// Final adds a function that is guaranteed to be executed
// even if an error is returned.
// Final functions are executed FIFO.
func (g *Group) Final(f func()) {
	g.final = append(g.final, f)
}

// Exec runs all functions then defer functions, stops on the
// first that errored and returns the error occurred.
// If no error is encountered, returns nil.
// If an error is returned, defer functions preceding the error
// returning function are executed.
func (g Group) Exec() error {
	defer func() {
		for _, f := range g.final {
			f()
		}
	}()

	var defers []func() error
	var err error
	for _, fn := range g.funcs {
		if fn.d {
			defers = append([]func() error{fn.f}, defers...)
			continue
		}
		if err = fn.f(); err != nil {
			break
		}
	}

	for _, fn := range defers {
		fn()
	}

	return err
}
