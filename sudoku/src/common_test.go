// these tests are verbose and can be simplified
// via iterator closures a la common.validate()

package common

import (
	"common"
	"testing"
)

var sampleCompleted [][]int = [][]int{
	{1, 2, 3, 4, 5, 6, 7, 8, 9},
	{4, 5, 6, 7, 8, 9, 1, 2, 3},
	{7, 8, 9, 1, 2, 3, 4, 5, 6},
	{2, 3, 4, 5, 6, 7, 8, 9, 1},
	{5, 6, 7, 8, 9, 1, 2, 3, 4},
	{8, 9, 1, 2, 3, 4, 5, 6, 7},
	{3, 4, 5, 6, 7, 8, 9, 1, 2},
	{6, 7, 8, 9, 1, 2, 3, 4, 5},
	{9, 1, 2, 3, 4, 5, 6, 7, 8}}

func TestIsValidRow(t *testing.T) {
	sudoku, _ := common.MakeSudoku(sampleCompleted)
	var res, ok bool
	var temp int

	// test valid row
	for r := 0; r < 9; r++ {
		res, ok = common.IsValidRow(sudoku, r)
		if !ok {
			t.Fatalf("bad input detected on row %d", r)
		}
		if !res {
			t.Fatalf("expected true on row %d", r)
		}

	}

	// test invalid row with bad input
	temp = sudoku.Values[0][0]
	sudoku.Values[0][0] = 10
	_, ok = common.IsValidRow(sudoku, 0)

	if ok {
		t.Fatalf("invalid values are not ok")
	}

	// test invalid row with duplicate value
	sudoku.Values[0][0] = temp
	sudoku.Values[0][0] = sudoku.Values[0][1]
	res, ok = common.IsValidRow(sudoku, 0)

	if res {
		t.Fatalf("duplicate values are not valid")
	}
}

func TestIsValidColumn(t *testing.T) {
	sudoku, _ := common.MakeSudoku(sampleCompleted)
	var res, ok bool
	var temp int

	// test valid column
	for c := 0; c < 9; c++ {
		res, ok = common.IsValidColumn(sudoku, c)
		if !ok {
			t.Fatalf("bad input detected on column %d", c)
		}
		if !res {
			t.Fatalf("expected true on column %d", c)
		}

	}

	// test invalid column with bad input
	temp = sudoku.Values[0][0]
	sudoku.Values[0][0] = 10
	_, ok = common.IsValidColumn(sudoku, 0)

	if ok {
		t.Fatalf("invalid values are not ok")
	}

	// test invalid column with duplicate value
	sudoku.Values[0][0] = temp
	sudoku.Values[0][0] = sudoku.Values[0][1]
	res, ok = common.IsValidColumn(sudoku, 0)

	if res {
		t.Fatalf("duplicate values are not valid")
	}
}

func TestIsValidBlock(t *testing.T) {
	sudoku, _ := common.MakeSudoku(sampleCompleted)
	var res, ok bool
	var temp int

	// test valid block
	for r := 0; r < 9; r++ {
		for  c := 0; c < 9; c++ {
			res, ok = common.IsValidBlock(sudoku, r, c)
			if !ok {
				t.Fatalf("bad input detected with row %d, column %d", r, c)
			}
			if !res {
				t.Fatalf("expected true with row %d, column %d", r, c)
			}
		}
	}

	// test invalid block with bad input
	temp = sudoku.Values[0][0]
	sudoku.Values[0][0] = 10
	_, ok = common.IsValidBlock(sudoku, 0, 0)

	if ok {
		t.Fatalf("invalid values are not ok")
	}

	// test invalid block with duplicate value
	sudoku.Values[0][0] = temp
	sudoku.Values[0][0] = sudoku.Values[0][1]
	res, ok = common.IsValidBlock(sudoku, 0, 0)

	if res {
		t.Fatalf("duplicate values are not valid")
	}
}
