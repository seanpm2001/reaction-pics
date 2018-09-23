package tumblr

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

// GetPosts returns a list of all Posts
func GetPosts(posts chan<- Post) {
	defer func() { close(posts) }()
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
