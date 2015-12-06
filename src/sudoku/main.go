package main

import (
	"fmt"
	"sudoku/examples"
	"sudoku/common"
	"sudoku/compute"
)

func printNdoku(ndoku *common.Ndoku) {
	var size int = int(ndoku.Size)
	var dim = size * size

	for r := 0; r < dim; r++ {
		for c := 0; c < dim; c++ {
			val := ndoku.Values[r][c]
			if val == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf("%d", ndoku.Values[r][c])
			}


			if c < dim-1 {
				fmt.Printf(" ")
			}
		}
		fmt.Println()
	}
}

func main() {
	ndoku, _ := common.MakeSudoku(examples.Ex2)
/*	ndoku.Values[0][0] = 0
	ndoku.Values[0][5] = 0
	ndoku.Values[5][5] = 0
	ndoku.Values[8][8] = 0
*/

	fmt.Println("solving the following problem:")
	printNdoku(ndoku) // TODO: change ndoku everywhere to sudoku for less goofiness
	res := compute.SolveDFS(ndoku)

	if res {
		fmt.Println()
		fmt.Println("solved:")
	} else {
		fmt.Println("could not solve: ")
	}

	printNdoku(ndoku)
}
