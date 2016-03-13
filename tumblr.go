package main

import _ "github.com/joho/godotenv/autoload"
import (
	"encoding/json"
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	"os"
)

func main() {
	posts := getPosts()
	fmt.Println(len(posts))
	fmt.Println(posts[0].Title)
	fmt.Println(posts[0].Body)
	fmt.Println(posts[0].Post_url)
}

func getPosts() []gotumblr.TextPost {
	client := gotumblr.NewTumblrRestClient(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("TOKEN"),
		os.Getenv("TOKEN_SECRET"),
		"https://www.albertyw.com/",
		"http://api.tumblr.com",
	)
	blogname := "devopsreactions.tumblr.com"
	blogtypes := "text"
	options := map[string]string{}
	postsResponse := client.Posts(blogname, blogtypes, options)
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
