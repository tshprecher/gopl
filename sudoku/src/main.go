package main

import (
	"fmt"
	"examples" // TODO: put these packages under sudoku/examples
	"common"
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
				fmt.Printf(", ")
			}
		}
		fmt.Println()
	}
}

func main() {
	fmt.Println("Ex 1:")
	ndoku, _ := common.MakeSudoku(examples.Ex1)
	printNdoku(ndoku)

	fmt.Println("Ex 2:")
	ndoku, _ = common.MakeSudoku(examples.Ex2)
	printNdoku(ndoku)
}
