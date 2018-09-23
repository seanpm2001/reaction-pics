package tumblr

import (
	"encoding/csv"
	"os"
	"strconv"
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

func getRow(post Post) (row []string) {
	row = []string{
		strconv.FormatInt(post.ID, 10),
		post.Title,
		post.URL,
		post.Image,
		strconv.FormatInt(post.Likes, 10),
	}
	return
}

func getCSVPath(blogName string) string {
	rootDir := os.Getenv("ROOT_DIR")
	return rootDir + "/" + csvDirectory + blogName + ".csv"
}
