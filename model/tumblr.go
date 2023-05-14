// Package model contains data for reaction.pics
package model

// GetPosts returns a list of all Posts
func GetPosts(posts chan<- Post) {
	defer func() { close(posts) }()
	csv := getCSV(false)
	existingPosts := ReadPostsFromCSV(csv)
	for _, p := range existingPosts {
		posts <- p
	}
}
