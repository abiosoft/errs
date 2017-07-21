package errs

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

var errFunc = func() error { return errors.New("error") }
var okFunc = func() error { return nil }

func TestNonNil(t *testing.T) {
	var e Group
	e.Add(errFunc)
	e.Add(okFunc)
	e.Add(okFunc)
	if e.Exec() == nil {
		t.Error("Expected error, found nil")
	}

	e = Group{}
	e.Add(okFunc)
	e.Add(errFunc)
	e.Add(okFunc)
	if e.Exec() == nil {
		t.Error("Expected error, found nil")
	}
}

func TestNil(t *testing.T) {
	var e Group
	e.Add(okFunc)
	e.Add(okFunc)
	e.Add(okFunc)
	if e.Exec() != nil {
		t.Error("Expected nil, found error")
	}

	e = Group{}
	if e.Exec() != nil {
		t.Error("Expected nil, found error")
	}
}

func TestDefer(t *testing.T) {
	var e Group
	var l []int
	e.Defer(func() {
		l = append(l, 3)
	})
	e.Add(func() error {
		l = append(l, 1)
		return nil
	})
	e.Defer(func() {
		l = append(l, 2)
	})
	if err := e.Exec(); err != nil {
		t.Errorf("Expected nil, found error: %v", err)
	}
	expected := []int{1, 2, 3}
	if fmt.Sprint(l) != fmt.Sprint(expected) {
		t.Errorf("Expected %v, found %v", expected, l)
	}
}

func TestFinal(t *testing.T) {
	var e Group
	e.Add(errFunc)
	var a, b int
	e.Final(func() { a = 100 })
	e.Final(func() { b = 101 })
	if e.Exec() == nil {
		t.Error("expected error, found nil")
	}
	if a != 100 {
		t.Errorf("Expected 100, found %v", a)
	}
	if b != 101 {
		t.Errorf("Expected 101, found %v", b)
	}
}

func TestFiller_Fill(t *testing.T) {
	var f = sliceFiller{
		0, "str", []int{2, 1},
	}
	var a int
	var b string
	var c []int
	f.Fill(&a, &b, &c)
	assert(t, a, f[0])
	assert(t, b, f[1])
	assert(t, c, f[2])
}

func TestGroup_AddF(t *testing.T) {
	var e Group
	var a, a1 int
	var b, b1 string
	a1 = 2
	b1 = "some string"
	var f = func(a int, b string) (int, string, error) {
		return a, b, nil
	}
	e.AddF(f, a1, b1).Fill(&a, &b)

	err := errors.New("err")
	var f1a int
	var f1a1 = 1
	var f1 = func() (int, error) { return f1a1, err }
	e.AddF(f1).Fill(&f1a)

	err1 := e.Exec()
	assert(t, err1, err)
	assert(t, a, a1)
	assert(t, b, b1)
	assert(t, f1a, f1a1)
}

func assert(t *testing.T, a, b interface{}) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%v != %v", a, b)
		t.Fail()
	}
}
