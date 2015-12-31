package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sudoku/common"
	"sudoku/compute"
	"sudoku/io"
	"time"
)

var stdoutWriter = io.NewWriter(os.Stdout)

func exitError(message string) {
	fmt.Println(fmt.Sprintf("error: %s", message))
	os.Exit(1)
}

func fetchSudoku(level, id int) *common.Sudoku {
	sudoku := io.FetchWebSudoku(level, id)
	if sudoku == nil {
		exitError("unexpected error fetching puzzle.")
	}
	return sudoku
}

func writeWebSudoku(sudoku *common.Sudoku, level, id int) {
	stdoutWriter.WriteComment("fetched from http://www.websudoku.com")
	stdoutWriter.WriteComment(fmt.Sprintf("level: %d, id: %d", level, id))
	stdoutWriter.WriteComment(fmt.Sprintf("%v", time.Now()))
	stdoutWriter.WriteSudoku(sudoku)
	fmt.Println()
}

func handleSolve(level, id int) {
	if level < 1 || level > 4 {
		exitError("arg 'level' must exist and be 1, 2, 3, or 4.")
	}
	if id < 0 {
		exitError("arg 'id' must exist and be positive.")
	}

	sudoku := fetchSudoku(level, id)
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

func handleDownload(n int) {
	if n == 0 {
		// read identifiers from stdin
		// NOTE: not concurrent => slower than it should be
		var level, id int
		read, err := fmt.Scanln(&level, &id);

		for read > 0 && err == nil {
			sudoku := fetchSudoku(level, id)
			writeWebSudoku(sudoku, level, id)
			read, err = fmt.Scanln(&level, &id)
		}
	} else {
		// download random identifiers
		if n < 0 {
			exitError("arg 'n' must be positive.")
		}
		for i := 0; i < n; i++ {
			level := rand.Intn(4) + 1
			id := rand.Intn(2e9)
			sudoku := fetchSudoku(level, id)
			writeWebSudoku(sudoku, level, id)
		}
	}

}

func main() {
	// for downloading
	dl := flag.Bool("dl", false, "download sudokus")
	n := flag.Int("n", 0, "download 'n' number of random puzzles")

	// for solving
	level := flag.Int("level", 0, "solve puzzle of difficulty 'level' (1, 2, 3, or 4)")
	id := flag.Int("id", -1, "solve puzzle with id 'id'")

	flag.Parse()

	if *dl {
		handleDownload(*n)
	} else {
		handleSolve(*level, *id)
	}
}
