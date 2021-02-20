// Package common contains logic for creating and validating puzzles.
package common

const EmptyField = 0

// A Sudoku represents a square sudoku puzzle of
// arbitrary size. A standard 9x9 sudoku has size == 3.
type Sudoku struct {
	Size   uint8
	Values [][]int
}

// NewSudoku creates a new Sudoku with a given size.
// The size is the sqrt of the length of one row.
func NewSudoku(size uint8) *Sudoku {
	dim := int(size) * int(size)
	values := make([][]int, dim)

	for d := 0; d < dim; d++ {
		values[d] = make([]int, dim)
	}

	return &Sudoku{size, values}
}

// NewSudokuFromSlice creates a new Sudoku given a two
// dimensional slice of record values and a size. On an error,
// ok == false and no puzzle is created.
func NewSudokuFromSlice(values [][]int, size uint8) (sudoku *Sudoku, ok bool) {
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

// IsValidRow returns true if a row's values
// do not violate the sudoku property.
func IsValidRow(sudoku *Sudoku, row int) bool {
	return validate(sudoku, row, 0, func(r int, c int, s uint8) (int, int, bool) {
		if c+1 >= int(s)*int(s) {
			return 0, 0, true
		}
		return r, c + 1, false
	})
}

// IsValidColumn returns true if a column's values
// do not violate the sudoku property.
func IsValidColumn(sudoku *Sudoku, col int) bool {
	return validate(sudoku, 0, col, func(r int, c int, s uint8) (int, int, bool) {
		if r+1 >= int(s)*int(s) {
			return 0, 0, true
		}
		return r + 1, c, false
	})
}

// IsValidBlock returns true if a block's values do not
// violate the sudoku property. The block is determined
// from the row and column indices.
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

// IsValid returns true if a Sudoku does not validate the sudoku properties.
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
