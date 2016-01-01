package compute

import (
	"sudoku/common"
	"sudoku/io"
	"sync"
)

type result struct {
	sudoku  *common.Sudoku
	success bool
}

func writeResults(writer *io.Writer, results []result) {
	for _, r := range results {
		writer.WriteComment("")
		if r.success {
			writer.WriteComment("SUCCESS")
		} else {
			writer.WriteComment("FAILURE")
		}
		writer.WriteComment("")
		writer.WriteSudoku(r.sudoku)
	}
}

func SolveSerial(algo func(*common.Sudoku) bool, puzzles []common.Sudoku, writer *io.Writer) {
	results := make([]result, len(puzzles))

	for p := range puzzles {
		res := algo(&puzzles[p])
		results[p] = result{&puzzles[p], res}
	}

	writeResults(writer, results)
}

func SolveParallel(algo func(*common.Sudoku) bool, puzzles []common.Sudoku, writer *io.Writer, concurrencyLevel uint8) {
	var wg sync.WaitGroup

	type message struct {
		index  int
		sudoku *common.Sudoku
	}

	ch := make(chan message, 25) // TODO: make this buffer size less arbitrary?
	results := make([]result, len(puzzles))

	for c := concurrencyLevel; c > 0; c-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for m := range ch {
				res := algo(m.sudoku)
				results[m.index] = result{m.sudoku, res} // TODO: do we need to synchronize here for many-core CPUs?
			}
		}()
	}

	for p := range puzzles {
		ch <- message{p, &puzzles[p]}
	}

	close(ch)

	wg.Wait()
	writeResults(writer, results)
}
