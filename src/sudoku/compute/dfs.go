/**
 * TODO: description, comments e'erywhere
 */
package compute

import (
	"sudoku/common"
)

func solveDFS(sudoku *common.Sudoku, curFieldNum int) bool {
	dim := int(sudoku.Size) * int(sudoku.Size)

	// no more fields to satisfy means problem is solved
	if curFieldNum == dim * dim {
		return true
	}

	curFieldRow, curFieldCol := curFieldNum / dim, curFieldNum % dim

	// if the field is already filled, move on to the next field
	if sudoku.Values[curFieldRow][curFieldCol] != common.EmptyField {
		return solveDFS(sudoku, curFieldNum+1)
	}

	// satisfy recursively with all field values
	for f := 1; f <= dim; f++ {
		sudoku.Values[curFieldRow][curFieldCol] = f

		// validate row, column, and block
		ok := common.IsValidRow(sudoku, curFieldRow)
		ok = ok && common.IsValidColumn(sudoku, curFieldCol)
		ok = ok && common.IsValidBlock(sudoku, curFieldRow, curFieldCol)

		if !ok {
			sudoku.Values[curFieldRow][curFieldCol] = common.EmptyField
			continue
		}

		ok = solveDFS(sudoku, curFieldNum+1)
		if ok {
			return true
		}
	}

	sudoku.Values[curFieldRow][curFieldCol] = common.EmptyField
	return false
}

func SolveDFS(sudoku *common.Sudoku) bool {
	ok := common.IsValid(sudoku)

	if !ok {
		return false
	}

	return solveDFS(sudoku, 0)
}
