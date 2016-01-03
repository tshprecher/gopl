package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sudoku/common"
	"sudoku/compute"
	"sudoku/io"
	"time"
)

var stdoutWriter = io.NewWriter(os.Stdout)

func fetchSudoku(level, id int) *common.Sudoku {
	sudoku := io.FetchWebSudoku(level, id)
	if sudoku == nil {
		log.Fatal("unexpected error fetching puzzle.")
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

func handleSolve(level, id int, par bool) {
	puzzles := make([]common.Sudoku, 0, 10)
	if level == 0 && id == -1 {
		// read puzzles from stdin and solve
		stdinReader := io.NewReader(os.Stdin)
		sudoku, err := stdinReader.Next()

		for sudoku != nil && err == nil {
			puzzles = append(puzzles, *sudoku)
			sudoku, err = stdinReader.Next()
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		if par {
			compute.SolveParallel(compute.SolveDFS, puzzles, stdoutWriter, 4)
		} else {
			compute.SolveSerial(compute.SolveDFS, puzzles, stdoutWriter)
		}

	} else {
		// solve single puzzle
		if level < 1 || level > 4 {
			log.Fatal("arg 'level' must exist and be 1, 2, 3, or 4.")
		}
		if id < 0 {
			log.Fatal("arg 'id' must exist and be positive.")
		}
		sudoku := fetchSudoku(level, id)
		puzzle := []common.Sudoku{*sudoku}
		compute.SolveSerial(compute.SolveDFS, puzzle, stdoutWriter)
	}
}

func handleDownload(n int) {
	if n == 0 {
		// read identifiers from stdin
		// NOTE: not concurrent => slower than it should be
		var level, id int
		read, err := fmt.Scanln(&level, &id)

		for read > 0 && err == nil {
			sudoku := fetchSudoku(level, id)
			writeWebSudoku(sudoku, level, id)
			read, err = fmt.Scanln(&level, &id)
		}
	} else {
		// download random identifiers
		if n < 0 {
			log.Fatal("arg 'n' must be positive.")
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
	runtime.GOMAXPROCS(runtime.NumCPU())
	// for downloading
	dl := flag.Bool("dl", false, "download sudokus")
	n := flag.Int("n", 0, "download 'n' number of random puzzles")

	// for solving
	level := flag.Int("level", 0, "solve puzzle of difficulty 'level' (1, 2, 3, or 4)")
	id := flag.Int("id", -1, "solve puzzle with id 'id'")
	par := flag.Bool("par", false, "solve puzzles in parallel")

	flag.Parse()

	if *dl {
		handleDownload(*n)
	} else {
		handleSolve(*level, *id, *par)
	}
}
