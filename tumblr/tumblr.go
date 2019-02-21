package tumblr

// GetPosts returns a list of all Posts
func GetPosts(posts chan<- Post) {
	defer func() { close(posts) }()
	csvPath := getCSVPath(false)
	existingPosts := ReadPostsFromCSV(csvPath)
	for _, p := range existingPosts {
		posts <- p
	}
}
