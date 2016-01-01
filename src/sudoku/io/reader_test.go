package io

import (
	"bytes"
	"testing"
)

func TestReadSingle(t *testing.T) {
	var tests = []struct {
		input string
		expect bool
	}{
		{"size 2 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4", true},
		{"size 2 . . . . 1 2 3 4 . . . . 1 2 3 4", true},
		{"size 2 1 2 3 4", false},
		{"size", false},
		{"size 1 1", true},
		{"size 1 .", true},
		{"size 1\n\n 1", true},
		{"size 1 /* block comment */ 1", true},
		{"size 1 \n// line comment\n1", true},
	}

	for _, test := range tests {
		reader := NewReader(bytes.NewBufferString(test.input))
		sud, err := reader.Read()
		if err != nil {
			if sud != nil {
				t.Errorf("unexpected sudoku value for input '%v'.", test.input)
			}
			if test.expect == true {
				t.Errorf("expected true, received err for input '%s': '%v'.", test.input, err)
			}
		} else {
			if sud == nil {
				t.Errorf("expected nil sudoku value for test '%v'.", test.input)
			}
			if test.expect == false {
				t.Errorf("expected false, received true for test '%v'.", test.input)
			}
		}
	}
}

func TestReadMultiple(t *testing.T) {
	success := "size 1 . size 1 1"
	reader := NewReader(bytes.NewBufferString(success))

	sud1, err1 := reader.Read()
	sud2, err2 := reader.Read()
	sud3, _ := reader.Read()

	if err1 != nil || err2 != nil {
		t.Errorf("unexpected error returned when reading multiple values.")
	}
	if sud1 == nil || sud2 == nil {
		t.Errorf("unexpected nil returned when reading multiple values.")
	}
	if sud3 != nil {
		t.Errorf("unexpected puzzle returned when done reading multiple values.")
	}

	fail := "size 1 . size 2 1"
	reader = NewReader(bytes.NewBufferString(fail))

	reader.Read()
	sud2, err2 = reader.Read()

	if err2 == nil {
		t.Errorf("expected an error when reading second puzzle for input '%v'.", fail)
	}
}
