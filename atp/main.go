package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strconv"
	"strings"
)

const (
	playerGetEndpoint = "http://www.tennis.com/rankings/ATP/"
	apiError          = "{\"error\": \"%s\"}\n"
)

// A player represents an ATP tennis player.
type player struct {
	Name    string `json:"name"`
	Rank    int    `json: "rank"`
	Country string `json: "rank"`
}

// writeErrorResponse writes the error message to the reponse in a standard form.
func writeErrorResponse(w http.ResponseWriter, status int, message string) {
	w.WriteHeader(status)
	fmt.Fprintf(w, apiError, message)
}

// filterHtmlNodes recursively filters html nodes given a filter function.
func filterHtmlNodes(node *html.Node, results *[]*html.Node, filter func(n *html.Node) bool) {
	if node == nil {
		return
	}
	if filter(node) {
		*results = append(*results, node)
	}
	child := node.FirstChild
	for child != nil {
		filterHtmlNodes(child, results, filter)
		child = child.NextSibling
	}
}

// handlePlayerGet serves the http endpoint to retrieve player information.
func handlePlayerGet(w http.ResponseWriter, r *http.Request) {
	lastName := strings.TrimSpace(strings.ToLower(r.URL.Query().Get("last_name")))
	if lastName == "" {
		writeErrorResponse(w, http.StatusBadRequest, "invalid last_name parameter")
		return
	}

	resp, err := http.Get(playerGetEndpoint)
	if err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	defer resp.Body.Close()

	var plr player
	p, _ := html.Parse(resp.Body)
	trNodes := make([]*html.Node, 0, 16)

	// get all the tr nodes represent each player's info
	filterHtmlNodes(p, &trNodes, func(n *html.Node) bool {
		return n.Type == html.ElementNode &&
			n.Data == "tr" &&
			len(n.Attr) > 0 &&
			(n.Attr[0].Val == "odd" || n.Attr[0].Val == "even")
	})

	// build the player for each record and compare against the last name given.
	for _, rec := range trNodes {
		curRecord := rec.FirstChild.NextSibling
		plr.Rank, _ = strconv.Atoi(curRecord.FirstChild.Data)
		curRecord = curRecord.NextSibling.NextSibling.NextSibling.NextSibling

		// the name may either be a raw string or inside an <a> tag
		if curRecord.FirstChild.NextSibling == nil {
			plr.Name = strings.TrimSpace(strings.ToLower(curRecord.FirstChild.Data))
		} else {
			plr.Name = strings.TrimSpace(strings.ToLower(curRecord.FirstChild.NextSibling.FirstChild.Data))
		}

		curRecord = curRecord.NextSibling.NextSibling
		plr.Country = strings.TrimSpace(strings.ToLower(curRecord.FirstChild.Data))

		if strings.Contains(plr.Name, lastName) {
			found, _ := json.Marshal(plr)
			fmt.Fprintf(w, string(found))
			fmt.Fprintln(w, "")
			return
		}
	}
	writeErrorResponse(w, http.StatusNotFound, "player not found")
}

func main() {
	http.HandleFunc("/player/get", handlePlayerGet)
	http.ListenAndServe(":8080", nil)
}
