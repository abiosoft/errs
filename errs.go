// Package errs is a convenience wrapper for chaining multiple
// error returning functions.
//  var e errs.Group
//
//  // add couple of functions
//  e.Add(func() error { ... })
//  e.Add(func() error { ... })
//  e.Add(func() error { ... })
//
//  // execute them
//  if err := e.Exec(); err != nil {
//      // handle error
//  }
package errs

// Group is a group of functions that returns an error.
// Empty value of Group is usable.
type Group []func() error

// Add adds f to the group of error functions.
func (g *Group) Add(f func() error) {
	*g = append(*g, f)
}

// Exec runs all functions serially and returns
// the first error occurred.
// If no error is encountered, returns nil.
func (g Group) Exec() error {
	for _, f := range g {
		if err := f(); err != nil {
			return err
		}
	}
	return nil
}
