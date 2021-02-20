// Package io implements logic for downloading puzzles from
// websudoku.com and writing/reading puzzles to/from files.
package io

import (
	"errors"
	"fmt"
	"github.com/tshprecher/gopl/sudoku/common"
	"io"
	"strconv"
	"text/scanner"
)

// A Reader deserializes an arbitrary number of puzzles
// from an input stream.
type Reader struct {
	// parsing state
	scanner      *scanner.Scanner
	curSudoku    *common.Sudoku
	curSudokuRow int

	// terminal flags
	termErr error
	termOk  bool
}

// NewReader creates a new Reader given an input stream.
func NewReader(reader io.Reader) *Reader {
	return &Reader{scanner: new(scanner.Scanner).Init(reader)}
}

// Next deserializes and returns the next puzzle in the stream.
// On success, a Sudoku is returned without an error. On failure,
// an error is returned without a Sudoku and all subsequent calls
// return the same error. If the Reader has reached the end of the
// stream, Next returns nil, nil.
func (r *Reader) Next() (*common.Sudoku, error) {
	// once a terminal state is reached, subsequent calls return the terminal state
	if r.termOk {
		return nil, nil
	}
	if r.termErr != nil {
		return nil, r.termErr
	}

	// attempt to read the next puzzle
	r.curSudoku = nil
	r.curSudokuRow = 0

	eof, size, err := r.scanSize()
	if eof {
		r.termOk = true
		return nil, nil
	}
	if err != nil {
		r.termErr = err
		return nil, err
	}

	r.curSudoku = common.NewSudoku(uint8(size))
	for s := r.curSudoku.Size * r.curSudoku.Size; s > 0; s-- {
		eof, err := r.scanRow()
		if eof {
			r.termOk = true
			return nil, nil
		}
		if err != nil {
			r.termErr = err
			return nil, err
		}
	}
	return r.curSudoku, nil
}

func (r *Reader) readError(message string) error {
	ln, col := r.scanner.Pos().Line, r.scanner.Pos().Column
	return errors.New(fmt.Sprintf("%s @ line %d, column %d.", message, ln, col))
}

func (r *Reader) scanSize() (eof bool, size int, err error) {
	eof, err = r.scanString("size")
	if eof {
		return true, 0, nil
	}
	if err != nil {
		return false, 0, err
	}
	return r.scanField()
}

func (r *Reader) scanRow() (eof bool, err error) {
	// expect r.curSudoku != nil
	for e, len := 0, int(r.curSudoku.Size)*int(r.curSudoku.Size); e < len; e++ {
		var val int
		eof, val, err = r.scanField()
		if eof {
			return true, nil
		}
		if err != nil {
			return false, err
		}

		if val != common.EmptyField && (val < 1 || val > len) {
			return false, r.readError(fmt.Sprintf("invalid value %d", val))
		}
		r.curSudoku.Values[r.curSudokuRow][e] = val
	}
	r.curSudokuRow++
	return false, nil
}

func (r *Reader) scanString(s string) (eof bool, err error) {
	// expect r.curSudoku != nil
	tok := r.scanner.Scan()
	if tok == scanner.EOF {
		return true, nil
	}
	if tokenText := r.scanner.TokenText(); tok != scanner.Ident || tokenText != s {
		return false, r.readError(fmt.Sprintf("unexpected identifier '%s'", tokenText))
	}
	return false, nil
}

func (r *Reader) scanField() (eof bool, field int, err error) {
	// expect r.curSudoku != nil
	tok := r.scanner.Scan()
	if tok == scanner.EOF {
		return true, 0, nil
	}
	if tok != scanner.Int && tok != '.' {
		return false, 0, r.readError(fmt.Sprintf("element not found"))
	}
	if tok == '.' {
		return false, common.EmptyField, nil
	}

	i, err := strconv.Atoi(r.scanner.TokenText())
	if err != nil {
		return false, 0, err
	}
	return false, i, nil
}
