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
func ReadPostsFromCSV(csvPath string) (posts []Post) {
	file, err := os.Open(csvPath)
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

func getCSVPath(test bool) string {
	path := prodCSVPath
	if test {
		path = testCSVPath
	}
	path = os.Getenv("ROOT_DIR") + "/" + path
	return path
}
