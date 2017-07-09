package tumblr

import (
	"encoding/json"
	"fmt"
	"github.com/MariaTerzieva/gotumblr"
	// Used for getting tumblr env vars
	_ "github.com/joho/godotenv/autoload"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

const (
	tumblrURL  = "http://api.tumblr.com"
	blogTypes  = "text"
	postsLimit = 20
)

// Blogs is a list of all blogs to read from
var Blogs = []string{
	"devopsreactions.tumblr.com",
	"lifeofasoftwareengineer.tumblr.com",
	"dbareactions.tumblr.com",
	"securityreactions.tumblr.com",
}

func duplicateChan(in <-chan Post, out1, out2 chan<- Post, wg *sync.WaitGroup) {
	for p := range in {
		out1 <- p
		out2 <- p
	}
	wg.Done()
}

// GetPosts returns a list of all Posts
func GetPosts(getNewPosts bool, out1 chan<- Post) {
	var wg sync.WaitGroup
	channels := [](chan<- Post){out1}
	for _, blogName := range Blogs {
		posts := make(chan Post)
		out2 := make(chan Post)
		channels = append(channels, out2)
		wg.Add(1)
		go duplicateChan(posts, out1, out2, &wg)
		go getBlogPosts(blogName, getNewPosts, posts)
		go WritePostsToCSV(blogName, out2)
	}
	wg.Wait()
	for _, c := range channels {
		close(c)
	}
}

func getBlogPosts(blogName string, getNewPosts bool, posts chan<- Post) {
	defer func() { close(posts) }()
	var newPosts []Post
	existingPosts := ReadPostsFromCSV(blogName)
	maxPostID := int64(0)
	for _, p := range existingPosts {
		if p.ID > maxPostID {
			maxPostID = p.ID
		}
		posts <- p
	}
	if !getNewPosts {
		return
	}
	offset := 0
	client := getTumblrClient()
	for len(newPosts) > 0 || offset == 0 {
		fmt.Println("Downloading", blogName, offset)
		options := getTumblrOptions(offset)
		postsResponse := client.Posts(blogName, blogTypes, options)
		newPosts = parsePosts(postsResponse)
		for _, p := range newPosts {
			if p.ID <= maxPostID {
				return
			}
			posts <- p
		}
		offset += postsLimit
	}
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
			err = getPostImage(post)
			if len(post.Title) == 0 {
				continue
			}
			if err == nil {
				posts = append(posts, *post)
			}
		}
	}
	return posts
}

func getPostImage(post *Post) error {
	imageName, imagePath := getImageNamePath(post.Image)
	output, err := os.Create(imagePath)
	if err != nil {
		fmt.Println("Cannot create ", imagePath)
		return err
	}
	defer output.Close()

	response, err := http.Get(post.Image)
	if err != nil {
		fmt.Println("Error downloading", post.Image, " - ", err)
		return err
	}
	defer response.Body.Close()

	_, err = io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error saving", post.Image)
		return err
	}

	post.Image = fmt.Sprintf(getImageURL(imageName))
	return nil
}

func getImageNamePath(imageName string) (name, path string) {
	name = uuid.NewV4().String() + filepath.Ext(imageName)
	rootDir := os.Getenv("ROOT_DIR")
	path = fmt.Sprintf("%s/tumblr/data/static/%s", rootDir, name)
	return
}

func getImageURL(imageName string) string {
	return fmt.Sprintf("/static/data/%s", imageName)
}
