// TODO: add tests

package io

import (
	"errors"
	"io"
	"fmt"
	"strconv"
	"sudoku/common"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(writer io.Writer) *Writer {
	return &Writer{writer}
}

func (w *Writer) WriteComment(message string) error {
	comment := fmt.Sprintf("//%s\n", message)
	_, err := w.writer.Write([]byte(comment))

	return err
}

func (w *Writer) WriteSudoku(sudoku *common.Sudoku) error {
	if sudoku == nil {
		return errors.New("cannot write nil sudoku.")
	}

	w.writer.Write([]byte(fmt.Sprintf("size %d\n", sudoku.Size)))
	for _, row := range sudoku.Values {
		for i, v := range row {
			// TODO: create once, resize as necessary?
			toWrite := make([]byte, 0, 8)

			if v == common.EmptyField {
				toWrite = append(toWrite, '.')
			} else {
				toWrite = append(toWrite, []byte(strconv.Itoa(v))...)
			}

			if i < len(row)-1 {
				toWrite = append(toWrite, ' ')
			} else {
				toWrite = append(toWrite, '\n')
			}

			if _, err := w.writer.Write(toWrite); err != nil {
				return err
			}
		}
	}
	return nil
}