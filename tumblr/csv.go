package tumblr

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const (
	prodCSVPath = "data/posts.csv"
	testCSVPath = "data/posts_test.csv"
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
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		filename = "."
	}
	path = filepath.Join(filepath.Dir(filename), path)
	fmt.Println(path)
	return path
}
