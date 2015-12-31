package main

import (
	"flag"
	"fmt"
	"os"
	"sudoku/compute"
	"sudoku/io"
)

var stdoutWriter = io.NewWriter(os.Stdout)

func exitError(message string) {
	fmt.Println(fmt.Sprintf("error: %s", message))
	os.Exit(1)
}

func handleSolve(level, id int) {
	if level < 1 || level > 4 {
		exitError("arg 'level' must exist and be 1, 2, 3, or 4.")
	}
	if id < 0 {
		exitError("arg 'id' must exist and be positive.")
	}

	sudoku := io.FetchWebSudoku(level, id)

	if sudoku == nil {
		exitError("problem not found")
	}

	fmt.Println("unsolved:")
	stdoutWriter.WriteSudoku(sudoku)

	res := compute.SolveDFS(sudoku)

	if res {
		fmt.Println()
		fmt.Println("solved:")
	} else {
		fmt.Println("could not solve: ")
	}

	stdoutWriter.WriteSudoku(sudoku)
}

func main() {
	download := flag.Bool("dl", false, "download sudoku(s)")

	solveLevel := flag.Int("level", 0, "problem level (1, 2, 3, 4)")
	solveId := flag.Int("id", -1, "problem id")

	flag.Parse()

	if *download {
		// TODO: implement the download case
		fmt.Println("error: download case not yet implemented")
	} else {
		handleSolve(*solveLevel, *solveId)
	}
}
