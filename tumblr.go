package main

import _ "github.com/joho/godotenv/autoload"
import (
	"encoding/json"
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	"os"
	"strconv"
)

const (
	TumblrURL  = "http://api.tumblr.com"
	BlogName   = "devopsreactions.tumblr.com"
	BlogTypes  = "text"
	PostsLimit = 20
)

func main() {
	posts := getPosts()
	fmt.Println(len(posts))
	fmt.Println(posts[0].Title)
	fmt.Println(posts[0].Body)
	fmt.Println(posts[0].Post_url)
}

func getPosts() []gotumblr.TextPost {
	var posts, newPosts []gotumblr.TextPost
	offset := 0
	client := getTumblrClient()
	for len(newPosts) == PostsLimit || offset == 0 {
		options := getTumblrOptions(offset)
		postsResponse := client.Posts(BlogName, BlogTypes, options)
		newPosts = parsePosts(postsResponse)
		posts = append(posts, newPosts...)
		offset += PostsLimit
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
		TumblrURL,
	)
	return client
}

func getTumblrOptions(offset int) map[string]string {
	options := map[string]string{}
	options["offset"] = strconv.Itoa(offset)
	options["limit"] = strconv.Itoa(PostsLimit)
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
