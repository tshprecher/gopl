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
	if ndoku.Values[curFieldRow][curFieldCol] != 0 {
		return solveDFS(ndoku, curFieldNum+1)
	}

	// satisfy recursively with all field values
	for f := 1; f <= dim; f++ {
		ndoku.Values[curFieldRow][curFieldCol] = f

		// validate row
		res, ok := common.IsValidRow(ndoku, curFieldRow)
		if !res || !ok {
			ndoku.Values[curFieldRow][curFieldCol] = 0
			continue
		}

		// validate column
		res, ok = common.IsValidColumn(ndoku, curFieldCol)
		if !res || !ok {
			ndoku.Values[curFieldRow][curFieldCol] = 0
			continue
		}

		// validate block
		res, ok = common.IsValidBlock(ndoku, curFieldRow, curFieldCol)
		if !res || !ok {
			ndoku.Values[curFieldRow][curFieldCol] = 0
			continue
		}

		res = solveDFS(ndoku, curFieldNum+1)
		if res {
			return true
		}
	}

	ndoku.Values[curFieldRow][curFieldCol] = 0
	return false
}

// TODO: check the input to see if it's valid to begin with
func SolveDFS(ndoku *common.Ndoku) bool {
	return solveDFS(ndoku, 0)
}
