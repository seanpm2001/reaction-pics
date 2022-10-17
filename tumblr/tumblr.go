// Package tumblr data model for reaction.pics
package tumblr

// GetPosts returns a list of all Posts
func GetPosts(posts chan<- Post) {
	defer func() { close(posts) }()
	csv := getCSV(false)
	existingPosts := ReadPostsFromCSV(csv)
	for _, p := range existingPosts {
		posts <- p
	}
}
