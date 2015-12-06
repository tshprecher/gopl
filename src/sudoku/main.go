package main

import (
	"fmt"
	"sudoku/examples"
	"sudoku/common"
	"sudoku/compute"
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
	sudoku, _ := common.MakeStandardSudoku(examples.Ex2)
/*	sudoku.Values[0][0] = 0
	sudoku.Values[0][5] = 0
	sudoku.Values[5][5] = 0
	sudoku.Values[8][8] = 0
*/

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
