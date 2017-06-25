package errs

import (
	"errors"
	"fmt"
	"testing"
)

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
	e.Defer(func() error {
		l = append(l, 3)
		return nil
	})
	e.Add(func() error {
		l = append(l, 1)
		return nil
	})
	e.Defer(func() error {
		l = append(l, 2)
		return nil
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
	var a int
	e.Final(func() { a = 100 })
	if e.Exec() == nil {
		t.Error("expected error, found nil")
	}
	if a != 100 {
		t.Errorf("Expected 100, found %v", a)
	}
}

var errFunc = func() error { return errors.New("error") }
var okFunc = func() error { return nil }
