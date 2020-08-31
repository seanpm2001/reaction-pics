package tumblr

import (
	"encoding/csv"
	"io"
	"os"
	"path/filepath"
)

const (
	prodCSVPath = "tumblr/data/posts.csv"
	testCSVPath = "tumblr/data/posts_test.csv"
)

// ReadPostsFromCSV reads a CSV file into a list of posts
func ReadPostsFromCSV(csvPath string) []Post {
	file, err := os.Open(csvPath)
	defer file.Close()
	if err != nil {
		return []Post{}
	}
	return readCSV(file)
}

func readCSV(data io.Reader) []Post {
	reader := csv.NewReader(data)
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

func getCSVPath(test bool) string {
	path := prodCSVPath
	if test {
		path = testCSVPath
	}
	path = filepath.Join(os.Getenv("ROOT_DIR"), path)
	return path
}
