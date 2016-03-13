package main

import _ "github.com/joho/godotenv/autoload"
import (
	"encoding/json"
	"fmt"
	"os"
	"github.com/MariaTerzieva/gotumblr"
)

func main() {

	client := gotumblr.NewTumblrRestClient(
		os.Getenv("CONSUMER_KEY"),
		os.Getenv("CONSUMER_SECRET"),
		os.Getenv("TOKEN"),
		os.Getenv("TOKEN_SECRET"),
		"https://www.albertyw.com/",
		"http://api.tumblr.com",
	)
	info := client.Info()
	fmt.Println(info.User.Name)
	blogname := "devopsreactions.tumblr.com"
	blogtypes := "text"
	options := map[string]string{}
	postsResponse := client.Posts(blogname, blogtypes, options)
	fmt.Println(postsResponse.Posts[0])
	var post gotumblr.BasePost
	err := json.Unmarshal(postsResponse.Posts[0], &post)
	fmt.Println(err)
	fmt.Println(post)
	//fmt.Println(postsResponse.Total_posts)
}
