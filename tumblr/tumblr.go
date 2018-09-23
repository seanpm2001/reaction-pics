package tumblr

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gofrs/uuid"
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
func GetPosts(posts chan<- Post) {
	defer func() { close(posts) }()
	for _, blogName := range Blogs {
		getBlogPosts(blogName, posts)
	}
}

func getBlogPosts(blogName string, posts chan<- Post) {
	existingPosts := ReadPostsFromCSV(blogName)
	maxPostID := int64(0)
	for _, p := range existingPosts {
		if p.ID > maxPostID {
			maxPostID = p.ID
		}
		posts <- p
	}
}

func getImageNamePath(imageName string) (name, path string, err error) {
	genUUID, err := uuid.NewV4()
	if err != nil {
		return
	}
	name = genUUID.String() + filepath.Ext(imageName)
	rootDir := os.Getenv("ROOT_DIR")
	path = fmt.Sprintf("%s/tumblr/data/static/%s", rootDir, name)
	return
}

func getImageURL(imageName string) string {
	return fmt.Sprintf("/static/data/%s", imageName)
}
