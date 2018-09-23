package tumblr

import (
	"encoding/csv"
	"os"
)

const csvDirectory = "tumblr/data/"

// ReadPostsFromCSV reads a CSV file into a list of posts
func ReadPostsFromCSV(blogName string) (posts []Post) {
	csvLocation := getCSVPath(blogName)
	file, err := os.Open(csvLocation)
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

func getCSVPath(blogName string) string {
	rootDir := os.Getenv("ROOT_DIR")
	return rootDir + "/" + csvDirectory + blogName + ".csv"
}
