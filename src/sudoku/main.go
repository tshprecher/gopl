package main

import (
	"fmt"
	"os"
	"strconv"
	"sudoku/common"
	"sudoku/compute"
	"sudoku/io"
)

func printSudoku(sudoku *common.Sudoku) {
	var size int = int(sudoku.Size)
	var dim = size * size

	for r := 0; r < dim; r++ {
		for c := 0; c < dim; c++ {
			val := sudoku.Values[r][c]
			if val == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%d", sudoku.Values[r][c])
			}


			if c < dim-1 {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}

func main() {
	// TODO: safer argument validation
	level, _ := strconv.Atoi(os.Args[1])
	id, _ := strconv.Atoi(os.Args[2])

	sudoku := io.FetchWebSudoku(level, id)

	if sudoku == nil {
		fmt.Println("problem not found")
		os.Exit(1)
	}

	fmt.Println("solving the following problem:")
	printSudoku(sudoku)
	res := compute.SolveDFS(sudoku)

	if res {
		fmt.Println()
		fmt.Println("solved:")
	} else {
		fmt.Println("could not solve: ")
	}

	printSudoku(sudoku)
}
