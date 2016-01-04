package io

import (
	"github.com/tshprecher/gopl/sudoku/common"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"strconv"
)

const (
	webSudokuEndpoint = "http://view.websudoku.com/"
)

func scrapeCheatAndEditMask(node *html.Node) (*string, *string) {
	if node == nil {
		return nil, nil
	}

	var cheat, editMask *string

	if node.Data == "input" {
		if len(node.Attr) >= 4 && node.Attr[0].Key == "name" && node.Attr[0].Val == "cheat" {
			cheat = &node.Attr[3].Val

		}
		if len(node.Attr) >= 3 && node.Attr[0].Key == "id" && node.Attr[0].Val == "editmask" {
			editMask = &node.Attr[2].Val
		}
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		c, e := scrapeCheatAndEditMask(c)
		if cheat == nil && c != nil {
			cheat = c
		}
		if editMask == nil && e != nil {
			editMask = e
		}

	}

	return cheat, editMask
}

// FetchWebSudoku scrapes a sudoku puzzle from websudoku.com.
// A Sudoku is returned on success, nil otherwise.
func FetchWebSudoku(level, id int) *common.Sudoku {
	resp, err := http.PostForm(webSudokuEndpoint,
		url.Values{"level": {strconv.Itoa(level)}, "set_id": {strconv.Itoa(id)}})

	if err != nil {
		return nil
	}

	defer resp.Body.Close()
	parsed, perr := html.Parse(resp.Body)

	if perr != nil {
		return nil
	}

	cheat, editMask := scrapeCheatAndEditMask(parsed)

	if cheat == nil || editMask == nil {
		return nil
	}

	if len(*cheat) != 81 || len(*editMask) != 81 {
		return nil
	}

	sudoku := common.NewSudoku(3)
	for i := 0; i < 81; i++ {
		if (*editMask)[i] == '0' {
			sudoku.Values[i/9][i%9] = int((*cheat)[i]) - int('0')
		}
	}

	return sudoku
}
