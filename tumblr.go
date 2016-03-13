package main

import _ "github.com/joho/godotenv/autoload"
import (
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
}
