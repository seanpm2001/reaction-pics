package tumblr

// GetPosts returns a list of all Posts
func GetPosts(posts chan<- Post) {
	defer func() { close(posts) }()
	existingPosts := ReadPostsFromCSV()
	for _, p := range existingPosts {
		posts <- p
	}
}
