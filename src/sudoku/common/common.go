package common

const EmptyField = 0

/*
  represent the general sudoku puzzle state
  where size == 3 => standard sudoku 9x9.
*/
type Sudoku struct {
	Size   uint8
	Values [][]int
}

func NewSudoku(size uint8) *Sudoku {
	dim := int(size) * int(size)
	values := make([][]int, dim)

	for d := 0; d < dim; d++ {
		values[d] = make([]int, dim)
	}

	return &Sudoku{size, values}
}

func NewSudokuFromSlice(values [][]int, size uint8) (*Sudoku, bool) {
	dim := int(size) * int(size)

	// verify dimensions of input and range of the values
	if len(values) != dim {
		return nil, false
	}
	for _, row := range values {
		if len(row) != dim {
			return nil, false
		}
		for _, v := range row {
			if v < 0 || v > dim {
				return nil, false
			}

		}
	}

	// copy the input size
	valuesCopy := make([][]int, len(values))
	for r, row := range values {
		rowCopy := make([]int, len(row))
		copy(rowCopy, row)
		valuesCopy[r] = rowCopy
	}
	return &Sudoku{size, valuesCopy}, true
}

func validate(sudoku *Sudoku, startRow int, startCol int, next func(row int, col int, size uint8) (int, int, bool)) bool {
	seen := make(map[int]bool)
	dim := int(sudoku.Size) * int(sudoku.Size)
	done := false

	for row, col := startRow, startCol; !done; row, col, done = next(row, col, sudoku.Size) {
		if row < 0 || row >= dim || col < 0 || col >= dim {
			return false
		}

		value := sudoku.Values[row][col]
		if value != EmptyField {
			if value < 1 || value > dim {
				return false
			}
			if seen[value] {
				return false
			}
			seen[value] = true
		}
	}
	return true
}

func IsValidRow(sudoku *Sudoku, row int) bool {
	return validate(sudoku, row, 0, func(r int, c int, s uint8) (int, int, bool) {
		if c+1 >= int(s)*int(s) {
			return 0, 0, true
		}
		return r, c + 1, false
	})
}

func IsValidColumn(sudoku *Sudoku, col int) bool {
	return validate(sudoku, 0, col, func(r int, c int, s uint8) (int, int, bool) {
		if r+1 >= int(s)*int(s) {
			return 0, 0, true
		}
		return r + 1, c, false
	})
}

func IsValidBlock(sudoku *Sudoku, row int, col int) bool {
	size := int(sudoku.Size)

	return validate(sudoku, row/size*size, col/size*size, func(r int, c int, s uint8) (int, int, bool) {
		size := int(s)

		if (c+1)/size == c/size {
			return r, c + 1, false
		} else {
			if (r+1)/size == r/size {
				return r + 1, c / size * size, false
			} else {
				return 0, 0, true

			}
		}
	})
}

func IsValid(sudoku *Sudoku) bool {
	dim := int(sudoku.Size) * int(sudoku.Size)

	for d := 0; d < dim; d++ {
		ok := IsValidRow(sudoku, d)
		if !ok {
			return false
		}

		ok = IsValidColumn(sudoku, d)
		if !ok {
			return false
		}
	}

	for br := 0; br < int(sudoku.Size); br++ {
		for bc := 0; bc < int(sudoku.Size); bc++ {
			ok := IsValidBlock(sudoku, br*int(sudoku.Size), bc*int(sudoku.Size))
			if !ok {
				return false
			}
		}
	}
	return true
}
