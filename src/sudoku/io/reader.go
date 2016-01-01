package io

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"sudoku/common"
	"text/scanner"
)

type readState struct {
	scanner      *scanner.Scanner
	curSudoku    *common.Sudoku
	curSudokuRow int
}

type Reader struct {
	state readState
}

func NewReader(reader io.Reader) *Reader {
	state := readState{scanner: new(scanner.Scanner).Init(reader)}
	return &Reader{state}
}

func (r *Reader) Read() (*common.Sudoku, error) {
	r.state.curSudoku = nil
	r.state.curSudokuRow = 0

	size, err := scanSize(&r.state)

	if err != nil {
		return nil, err
	}

	r.state.curSudoku = common.NewSudoku(uint8(size))
	for s := r.state.curSudoku.Size * r.state.curSudoku.Size; s > 0; s-- {
		if err := scanRow(&r.state); err != nil {
			return nil, err
		}
	}

	return r.state.curSudoku, nil
}

// TODO: fix the positioning since the position returned if after the field we just scanned?
func readError(state *readState, message string) error {
	return errors.New(fmt.Sprintf("%s @ line %d, column %d.", message, state.scanner.Pos().Line, state.scanner.Pos().Column))
}

func scanSize(state *readState) (int, error) {
	if err := scanString(state, "size"); err != nil {
		return 0, err
	}
	size, err := scanField(state);

	if err != nil {
		return 0, err

	}

	return size, nil
}

func scanRow(state *readState) error {
	// expect state.curSudoku != nil
	for e, len := 0, int(state.curSudoku.Size)*int(state.curSudoku.Size); e < len; e++ {
		{
			val, err := scanField(state)
			if err != nil {
				return err
			}

			if val != common.EmptyField && (val < 1 || val > len) {
				return readError(state, fmt.Sprintf("invalid value %d", val))
			}
			state.curSudoku.Values[state.curSudokuRow][e] = val
		}
	}
	state.curSudokuRow++
	return nil
}


func scanString(state *readState, s string) error {
	// expect state.curSudoku != nil
	if tok := state.scanner.Scan(); tok != scanner.Ident {
		return readError(state, "unexpected error")
	}

	if tokenText := state.scanner.TokenText(); tokenText != s {
		return readError(state, fmt.Sprintf("unexpected identifier '%s'", tokenText))
	}

	return nil
}

func scanField(state *readState) (int, error) {
	// expect state.curSudoku != nil
	tok := state.scanner.Scan();

	if tok != scanner.Int && tok != '.' {
		return 0, readError(state, fmt.Sprintf("element not found"))
	}
	if tok == '.' {
		return common.EmptyField, nil
	}

	i, err := strconv.Atoi(state.scanner.TokenText())
	if err != nil {
		return 0, err
	}
	return i, nil
}
