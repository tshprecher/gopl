package main

import (
	"fmt"
	"examples"
)

func printProblem(problem [][]int) {
	for _, row := range(problem) {
		fmt.Println("", row)
	}
}

func main() {
	fmt.Println("testing compile")
	printProblem(examples.Ex1)
}
