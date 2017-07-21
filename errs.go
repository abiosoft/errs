// Package errs is a convenience wrapper for chaining multiple
// error returning functions when you do not need to handle the
// errors separately.
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

import (
	"fmt"
	"reflect"
	"sync"
)

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

// AddF is like Add but takes in a function
// and its arguments for cleaner code.
//
// Instead of the following
//  var n int64
//  e.Add(func() (err error) {
//    n, err = io.Copy(dst, src)
//    return
//  })
// You can have the following
//  var n int64
//  e.AddF(io.Copy, dst, src).Fill(&n)
// Explanation: io.Copy will be called with dst and src passed as arguments.
// if f returns one or more types, all types apart from error can
// be retrieved in order using the returned Filler as shown above.
// The Filler waits for the function to execute before filling the
// values.
func (g *Group) AddF(f interface{}, args ...interface{}) Filler {
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		panic(fmt.Errorf("Cannot invoke non function type '%v'", t))
	}
	rArgs := make([]reflect.Value, t.NumIn())
	for i := range rArgs {
		rArgs[i] = reflect.ValueOf(args[i])
	}
	var filler asyncFiller
	g.Add(func() (err error) {
		vals := reflect.ValueOf(f).Call(rArgs)
		var out []interface{}
		if len(vals) > 0 {
			for i := range vals {
				if err1, ok := vals[i].Interface().(error); ok {
					err = err1
					continue
				}
				out = append(out, vals[i].Interface())
			}
		}
		filler.set(out)
		return
	})
	return &filler
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

// Filler fills in the return values of a
// function executed by Group.AddF into args.
type Filler interface {
	// Fill fills in the values into args.
	Fill(args ...interface{})
	// FillAt fills in the value at i into arg.
	FillAt(i int, arg interface{})
}

type sliceFiller []interface{}

func (f sliceFiller) Fill(args ...interface{}) {
	for i := 0; i < len(f) && i < len(args); i++ {
		f.FillAt(i, args[i])
	}
}
func (f sliceFiller) FillAt(i int, arg interface{}) {
	reflect.ValueOf(arg).Elem().Set(reflect.ValueOf(f[i]))
}

type asyncFiller struct {
	values []interface{}
	toFill []struct {
		i   int
		arg interface{}
	}
	ready bool
	sync.Mutex
}

func (f *asyncFiller) set(v []interface{}) {
	f.Lock()
	defer f.Unlock()
	f.values = v
	f.ready = true
	for i := range f.toFill {
		sliceFiller(f.values).FillAt(f.toFill[i].i, f.toFill[i].arg)
	}
	f.toFill = nil
}

func (f *asyncFiller) Fill(args ...interface{}) {
	for i := range args {
		f.FillAt(i, args[i])
	}
}
func (f *asyncFiller) FillAt(i int, arg interface{}) {
	f.Lock()
	defer f.Unlock()
	if f.ready {
		sliceFiller(f.values).FillAt(i, arg)
	} else {
		f.toFill = append(f.toFill, struct {
			i   int
			arg interface{}
		}{i: i, arg: arg})
	}
}
