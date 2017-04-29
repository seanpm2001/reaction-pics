package tumblr

import (
	"encoding/json"
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	// Used for getting tumblr env vars
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
)

const (
	tumblrURL  = "http://api.tumblr.com"
	blogName   = "devopsreactions.tumblr.com"
	blogTypes  = "text"
	postsLimit = 20
)

// GetPosts returns a list of all Posts
func GetPosts(getNewPosts bool) []Post {
	var newPosts []Post
	posts := ReadPostsFromCSV()
	if !getNewPosts {
		return posts
	}
	maxPostID := int64(0)
	for _, p := range posts {
		if p.ID > maxPostID {
			maxPostID = p.ID
		}
	}
	offset := 0
	client := getTumblrClient()
	for len(newPosts) == postsLimit || offset == 0 {
		fmt.Println("Downloading", offset)
		options := getTumblrOptions(offset)
		postsResponse := client.Posts(blogName, blogTypes, options)
		newPosts = parsePosts(postsResponse)
		newPosts[5].ID = 1
		for _, p := range newPosts {
			if p.ID <= maxPostID {
				return posts
			}
			posts = append(posts, p)
		}
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
	options["notes_info"] = strconv.FormatBool(true)
	return options
}

func parsePosts(postsResponse gotumblr.PostsResponse) []Post {
	var posts []Post
	var tumblrPost gotumblr.TextPost
	for _, element := range postsResponse.Posts {
		err := json.Unmarshal(element, &tumblrPost)
		if err != nil {
			fmt.Println(err)
		} else {
			post := GoTumblrToPost(&tumblrPost)
			posts = append(posts, *post)
		}
	}
	return posts
}
