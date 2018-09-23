package tumblr

import (
	"encoding/csv"
	"os"
)

const (
	prodCSVPath = "tumblr/data/posts.csv"
	testCSVPath = "tumblr/data/posts_test.csv"
)

// ReadPostsFromCSV reads a CSV file into a list of posts
func ReadPostsFromCSV() (posts []Post) {
	file, err := os.Open(getCSVPath())
	defer file.Close()
	if err != nil {
		return
	}
	reader := csv.NewReader(file)
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

func getCSVPath() string {
	env := os.Getenv("ENVIRONMENT")
	path := testCSVPath
	if env == "production" {
		path = prodCSVPath
	}
	path = os.Getenv("ROOT_DIR") + "/" + path
	return path
}
