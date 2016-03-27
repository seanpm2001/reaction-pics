package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/MariaTerzieva/gotumblr"
	"os"
	"strconv"
)

const (
	tumblrURL  = "http://api.tumblr.com"
	blogName   = "devopsreactions.tumblr.com"
	blogTypes  = "text"
	postsLimit = 20
)

func main() {
	posts := getPosts()
	writePostsToCSV(posts)
}

func writePostsToCSV(posts []gotumblr.TextPost) {
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
		post.Body,
		post.Post_url,
	}
	return row
}

func getPosts() []gotumblr.TextPost {
	var posts, newPosts []gotumblr.TextPost
	offset := 0
	client := getTumblrClient()
	for len(newPosts) == postsLimit || offset == 0 {
		options := getTumblrOptions(offset)
		postsResponse := client.Posts(blogName, blogTypes, options)
		newPosts = parsePosts(postsResponse)
		posts = append(posts, newPosts...)
		offset += postsLimit
	}
	return posts
}

func getTumblrClient() *gotumblr.TumblrRestClient {
	client := gotumblr.NewTumblrRestClient(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("TOKEN"),
		os.Getenv("TOKEN_SECRET"),
		"https://www.albertyw.com/",
		tumblrURL,
	)
	return client
}

func getTumblrOptions(offset int) map[string]string {
	options := map[string]string{}
	options["offset"] = strconv.Itoa(offset)
	options["limit"] = strconv.Itoa(postsLimit)
	return options
}

func parsePosts(postsResponse gotumblr.PostsResponse) []gotumblr.TextPost {
	var posts []gotumblr.TextPost
	var post gotumblr.TextPost
	for _, element := range postsResponse.Posts {
		err := json.Unmarshal(element, &post)
		if err != nil {
			fmt.Println(err)
		} else {
			posts = append(posts, post)
		}
	}
	return posts
}
