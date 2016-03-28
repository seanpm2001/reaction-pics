package tumblr

import (
	"encoding/csv"
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	"os"
	"strconv"
)

func WritePostsToCSV(posts []gotumblr.TextPost) {
	var row []string
	file, err := os.Create("data.csv")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	writer := csv.NewWriter(file)
	for _, post := range posts {
		row = getRow(post)
		writer.Write(row)
		if err := writer.Error(); err != nil {
			fmt.Println(err)
		}
	}
	writer.Flush()
}

func getRow(post gotumblr.TextPost) (row []string) {
	row = []string{
		strconv.FormatInt(post.Id, 10),
		post.Title,
		post.Post_url,
	}
	return
}
