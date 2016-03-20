package main

import _ "github.com/joho/godotenv/autoload"
import (
	"encoding/json"
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	"os"
)

const (
	TumblrURL = "http://api.tumblr.com"
	BlogName  = "devopsreactions.tumblr.com"
	BlogTypes = "text"
)

func main() {
	posts := getPosts()
	fmt.Println(len(posts))
	fmt.Println(posts[0].Title)
	fmt.Println(posts[0].Body)
	fmt.Println(posts[0].Post_url)
}

func getTumblrClient() *gotumblr.TumblrRestClient {
	client := gotumblr.NewTumblrRestClient(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("TOKEN"),
		os.Getenv("TOKEN_SECRET"),
		"https://www.albertyw.com/",
		TumblrURL,
	)
	return client
}

func getPosts() []gotumblr.TextPost {
	client := getTumblrClient()
	options := map[string]string{}
	options["offset"] = "1"
	options["limit"] = "5"
	postsResponse := client.Posts(BlogName, BlogTypes, options)
	posts := parsePosts(postsResponse)
	return posts
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
