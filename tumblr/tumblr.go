package tumblr

import (
	"encoding/json"
	"errors"
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
	posts, postIds := getExistingPosts()
	if !getNewPosts {
		return posts
	}
	offset := 0
	client := getTumblrClient()
	for len(newPosts) == postsLimit || offset == 0 {
		fmt.Println("Downloading", offset)
		options := getTumblrOptions(offset)
		postsResponse := client.Posts(blogName, blogTypes, options)
		newPosts = parsePosts(postsResponse)
		err := addNewPosts(&newPosts, &posts, &postIds)
		if err != nil {
			break
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

func getExistingPosts() ([]Post, map[int64]bool) {
	posts := ReadPostsFromCSV()
	postIds := make(map[int64]bool)
	for i := 0; i < len(posts); i++ {
		postIds[posts[i].ID] = true
	}
	return posts, postIds
}

func addNewPosts(newPosts *[]Post, posts *[]Post, postIds *map[int64]bool) (err error) {
	for i := 0; i < len(*newPosts); i++ {
		post := (*newPosts)[i]
		if _, ok := (*postIds)[post.ID]; ok {
			return errors.New("repeated data")
		}
		*posts = append(*posts, post)
	}
	return nil
}
