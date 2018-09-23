package tumblr

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

// Blogs is a list of all blogs to read from
var Blogs = []string{
	"devopsreactions.tumblr.com",
	"lifeofasoftwareengineer.tumblr.com",
	"dbareactions.tumblr.com",
	"securityreactions.tumblr.com",
}

// GetPosts returns a list of all Posts
func GetPosts(posts chan<- Post) {
	defer func() { close(posts) }()
	for range Blogs {
		getBlogPosts(posts)
	}
}

func getBlogPosts(posts chan<- Post) {
	existingPosts := ReadPostsFromCSV()
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
