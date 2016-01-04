package io

import (
	"bytes"
	"sudoku/common"
	"testing"
)

func TestNextSingle(t *testing.T) {
	reader := NewReader(bytes.NewBufferString("size 2 1 2 3 4 . . . . 1 2 3 4 . . . ."))
	sud, _ := reader.Next()

	if sud == nil || sud.Size != 2 {
		t.Errorf("expected sudoku with size = 2, received %v", sud)
	}

	for r := 0; r < 4; r++ {
		for c := 0; c < 4; c++ {
			if r % 2 == 0 && sud.Values[r][c] != c+1 {
				t.Errorf("expected value %d, received value %d", c+1, sud.Values[r][c])
			}
			if r % 2 == 1 && sud.Values[r][c] != common.EmptyField {
				t.Errorf("expected value %d, received value %d", common.EmptyField, sud.Values[r][c])
			}
		}
	}
}

func TestNextMultiple(t *testing.T) {
	var tests = []struct {
		input  string
		expect []bool // true => returns a sudoku
	}{
		{"size 2 1 2 3 4 1 2 3 4 1 2 3 4 1 2 3 4", []bool{true, false}},
		{"size 2 . . . . 1 2 3 4 . . . . 1 2 3 4", []bool{true, false}},
		{"size 2 1 2 3 4", []bool{false}},
		{"size 1 1 2 3 4", []bool{true, false}},
		{"size", []bool{false}},
		{"size 1 A", []bool{false}},
		{"size 1 1", []bool{true, false}},
		{"size 1 .", []bool{true, false}},
		{"size 1\n\n 1", []bool{true}},
		{"size 1 /* block comment */ 1", []bool{true, false}},
		{"size 1 \n// line comment\n1", []bool{true, false}},
	}

	for _, test := range tests {
		reader := NewReader(bytes.NewBufferString(test.input))
		for _, expect := range test.expect {
			sud, err := reader.Next()
			if err != nil {
				if sud != nil {
					t.Errorf("unexpected sudoku value along with err for input '%v'.", test.input)
				}
				if expect == true {
					t.Errorf("expected sudoku, received err for input '%s': '%v'.", test.input, err)
				}
			} else {
				if sud == nil && expect == true {
					t.Errorf("received EOF for test '%v'.", test.input)
				}
				if sud != nil && expect == false {
					t.Errorf("expected err, received sudoku for test '%v'.", test.input)
				}
			}
		}
	}
}

func TestTerminalError(t *testing.T) {
	reader := NewReader(bytes.NewBufferString("size 1 A"))
	_, err := reader.Next()
	_, err2 := reader.Next()

	if err == nil || err != err2 {
		t.Errorf("expected err == err2, received err = '%v', err2 = '%v'.", err, err2)
	}
}

func TestTerminalOk(t *testing.T) {
	reader := NewReader(bytes.NewBufferString("size 1 1"))
	// read the initial puzzle
	reader.Next()
	// should return nil, nil both times
	sud, err := reader.Next()
	sud2, err2 := reader.Next()

	if sud != nil || err != nil || sud2 != nil || err2 != nil {
		t.Errorf("expected sud == err == sud2 == err2 == nil")
	}
}
