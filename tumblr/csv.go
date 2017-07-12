package tumblr

import (
	"encoding/csv"
	"fmt"
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

// WritePostsToCSV writes a list of posts to a CSV file
func WritePostsToCSV(blogName string, postChan <-chan Post) (csvLocation string) {
	board := Board{}
	for p := range postChan {
		board.AddPost(p)
	}
	fmt.Printf("Saving %d posts\n", len(board.Posts))
	csvLocation = getCSVPath(blogName)
	file, err := os.Create(csvLocation)
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	board.SortPostsByID()
	writer := csv.NewWriter(file)
	for _, post := range board.Posts {
		row := getRow(post)
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
	return
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
