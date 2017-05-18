package errs

import (
	"errors"
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

var errFunc = func() error { return errors.New("error") }
var okFunc = func() error { return nil }
