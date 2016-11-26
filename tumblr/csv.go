package tumblr

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

const csvLocation = "data.csv"

// ReadPostsFromCSV reads a CSV file into a list of posts
func ReadPostsFromCSV() (posts []Post) {
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

// WritePostsToCSV writes a list of posts to a CSV file
func WritePostsToCSV(posts []Post) {
	var row []string
	fmt.Printf("Saving %d posts\n", len(posts))
	file, err := os.Create(csvLocation)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	posts = *SortPosts(&posts)
	writer := csv.NewWriter(file)
	for _, post := range posts {
		row = getRow(post)
		writer.Write(row)
		if err := writer.Error(); err != nil {
			fmt.Println(err)
		}
	}
	writer.Flush()
	err = writer.Error()
	if err != nil {
		fmt.Println(err)
	}
}

func getRow(post Post) (row []string) {
	row = []string{
		strconv.FormatInt(post.ID, 10),
		post.Title,
		post.URL,
	}
	return
}
