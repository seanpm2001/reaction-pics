package model

import (
	"bytes"
	_ "embed"
	"encoding/csv"
)

//go:embed data/posts.csv
var prodCSV []byte

//go:embed data/posts_test.csv
var testCSV []byte

// ReadPostsFromCSV reads a CSV file into a list of posts
func ReadPostsFromCSV(csvData []byte) []Post {
	ioReader := bytes.NewReader(csvData)
	reader := csv.NewReader(ioReader)
	var posts []Post
	for {
		row, err := reader.Read()
		if err != nil {
			break
		}
		post := CSVToPost(row)
		posts = append(posts, *post)
	}
	return posts
}

func getCSV(test bool) []byte {
	if test {
		return testCSV
	}
	return prodCSV
}
