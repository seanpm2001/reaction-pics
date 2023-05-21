package model

import (
	_ "embed"
	"strings"
)

// generated from https://gist.github.com/rg089/35e00abf8941d72d419224cfd5b5925d
//
//go:embed data/stopwords.txt
var stopwords []byte

// loadStopwords generates a map (set) of stopwords
func loadStopwords() map[string]bool {
	wordArray := strings.Split(string(stopwords), "\n")
	words := make(map[string]bool, len(wordArray))
	for _, word := range wordArray {
		words[word] = true
	}
	return words
}
