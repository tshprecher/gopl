package io

import (
	"bytes"
	"testing"
)

func TestReader(t *testing.T) {
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
				t.Errorf("unexpected sudoku value for input '%v'", test.input)
			}
			if test.expect == true {
				t.Errorf("expected true, received err for input '%s': '%v'", test.input, err)
			}
		} else {
			if sud == nil {
				t.Errorf("expected nil sudoku value for test '%v'", test.input)
			}
			if test.expect == false {
				t.Errorf("expected false, received true for test '%v'", test.input)
			}
		}
	}
}
