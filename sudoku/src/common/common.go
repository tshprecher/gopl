// common logic
package common

import "fmt"

const EMPTY_FIELD = 0

/*
  represent the general n-doku puzzle state
  where size == 3 => standard sudoku 9x9.
 */
type Ndoku struct{
	size uint8
	values [][]int
}

func MakeSudoku(values [][]int) (*Ndoku, bool) {
	return MakeNdoku(values, 3)
}

func MakeNdoku(values [][]int, size uint8) (*Ndoku, bool) {
	fmt.Println("making ndoku...")
	// verify dimensions of input and range of the values
	if len(values) != int(size) {
		return nil, false
	}
	for _, row := range(values) {
		if len(row) != int(size) {
			return nil, false
		}
		for _, v := range(row) {
			if v < 0 || v > int(size)*int(size) {
				return nil, false
			}

		}
	}

	// copy the input size
	valuesCopy := make([][]int, len(values))
	for r, row := range(values) {
		rowCopy := make([]int, len(row))
		copy(rowCopy, row)
		valuesCopy[r] = row
	}

	fmt.Println("made ndoku %v", &Ndoku{size, valuesCopy})

	return &Ndoku{size, valuesCopy}, true
}

func validate(ndoku *Ndoku, startRow int, startCol int, next func (row int, col int, size uint8) (int, int, bool)) (isValid, ok bool) {
	seen := make(map[int]bool)
	dim := int(ndoku.size)*int(ndoku.size)
	row, col, done := next(startRow, startCol, ndoku.size)

	for ; !done; row, col, done = next(row, col, ndoku.size) {
		if row < 0 || row >= dim || col < 0 || col >= dim {
			return false, false
		}

		value := ndoku.values[row][col]

		if value < 1 || value > dim {
			return false, false
		}
		if seen[value] {
			return false, true
		}
		seen[value] = true
	}
	return true, true
}

func IsValidRow(ndoku *Ndoku, row int) (bool, bool) {
	return validate(ndoku, row, 0, func (r int, c int, s uint8) (int, int, bool) {
		if c >= int(s)*int(s) {
			return 0, 0, true
		}
		return r, c+1, false
	})
}

func IsValidColumn(ndoku *Ndoku, col int) (bool, bool) {
	return validate(ndoku, 0, col, func (r int, c int, s uint8) (int, int, bool) {
		if r >= int(s)*int(s) {
			return 0, 0, true
		}
		return r+1, c, false
	})
}


func IsValidBlock(ndoku *Ndoku, row int, col int) (bool, bool) {
	startRow := row / int(ndoku.size)
	startCol := col / int(ndoku.size)
	count := 0

	return validate(ndoku, startRow, startCol, func (r int, c int, s uint8) (int, int, bool) {
		if count >= int(s)*int(s) {
			return 0, 0, true
		}
		count++
		return r+(count % 3), c + (count / 3), false
	})
}
