/**
 * TODO: description, comments e'erywhere
 */
package compute

import (
	"sudoku/common"
)

func solveDFS(ndoku *common.Ndoku, curFieldNum int) bool {
	dim := int(ndoku.Size) * int(ndoku.Size)

	// no more fields to satisfy means problem is solved
	if curFieldNum == dim * dim {
		return true
	}

	curFieldRow, curFieldCol := curFieldNum / dim, curFieldNum % dim

	// if the field is already filled, move on to the next field
	if ndoku.Values[curFieldRow][curFieldCol] != common.EmptyField {
		return solveDFS(ndoku, curFieldNum+1)
	}

	// satisfy recursively with all field values
	for f := 1; f <= dim; f++ {
		ndoku.Values[curFieldRow][curFieldCol] = f

		// validate row, column, and block
		ok := common.IsValidRow(ndoku, curFieldRow)
		ok = ok && common.IsValidColumn(ndoku, curFieldCol)
		ok = ok && common.IsValidBlock(ndoku, curFieldRow, curFieldCol)

		if !ok {
			ndoku.Values[curFieldRow][curFieldCol] = common.EmptyField
			continue
		}

		ok = solveDFS(ndoku, curFieldNum+1)
		if ok {
			return true
		}
	}

	ndoku.Values[curFieldRow][curFieldCol] = common.EmptyField
	return false
}

func SolveDFS(ndoku *common.Ndoku) bool {
	ok := common.IsValid(ndoku)

	if !ok {
		return false
	}

	return solveDFS(ndoku, 0)
}
