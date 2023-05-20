// Package model contains data for reaction.pics
package model

import (
	"encoding/json"
	"math/rand"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/gosimple/slug"
	"github.com/rollbar/rollbar-go"
)

// MaxKeywords is the maximum number of keywords that can be returned by a board
const (
	MaxKeywords   = 20
	imageRootPath = "https://img.reaction.pics/file/reaction-pics/"
)

// Post is a representation of a single tumblr post
type Post struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	URL   string `json:"url"`
	Image string `json:"image"`
	Likes int64  `json:"likes"`
}

// InternalURL returns the path to the post
func (p Post) InternalURL() string {
	slug := slug.Make(p.Title)
	if len(slug) > 30 {
		slug = slug[0:30]
	}
	return "/post/" + strconv.FormatInt(p.ID, 10) + "/" + slug
}

// MarshalJSON allows Post to be converted to json
func (p Post) MarshalJSON() ([]byte, error) {
	type jPost Post
	return json.Marshal(&struct {
		jPost
		InternalURL string `json:"internalURL"`
	}{
		jPost:       jPost(p),
		InternalURL: p.InternalURL(),
	})
}

// CSVToPost converts a CSV row into a Post
func CSVToPost(row []string) *Post {
	id, err := strconv.ParseInt(row[0], 10, 64)
	if err != nil {
		rollbar.Error("cannot parse id", map[string]string{"id": row[0]})
		id = 0
	}

	imageURL := imageRootPath + row[3]

	likes, err := strconv.ParseInt(row[4], 10, 64)
	if err != nil {
		rollbar.Error("cannot parse likes", map[string]string{"likes": row[4]})
		likes = 0
	}
	post := Post{
		ID:    id,
		Title: row[1],
		URL:   row[2],
		Image: imageURL,
		Likes: likes,
	}
	return &post
}

// Board is a container for Posts that offers serialization, sorting, and
// parallelization
type Board struct {
	Posts []Post
	mut   *sync.RWMutex
}

// InitializeBoard means to create a new board and start writing reading saved
// posts into it
func InitializeBoard() *Board {
	board := NewBoard([]Post{})
	board.PopulateBoard()
	return &board
}

// NewBoard creates a Board from an array of Posts
func NewBoard(p []Post) Board {
	return Board{
		Posts: p,
		mut:   &sync.RWMutex{},
	}
}

// PopulateBoard adds production posts to the board
func (b *Board) PopulateBoard() {
	b.mut.Lock()
	go func() {
		defer b.mut.Unlock()
		b.Posts = ReadPostsFromCSV(getCSV(false))
		b.sortPostsByLikes()
	}()
}

// AddPost adds a single post to the board; no-ops if the post is already present
func (b *Board) AddPost(p Post) {
	b.mut.Lock()
	defer b.mut.Unlock()
	for i := 0; i < len(b.Posts); i++ {
		if b.Posts[i].Title == p.Title {
			return
		}
	}
	b.Posts = append(b.Posts, p)
}

// MarshalJSON allows Board to be converted to json
func (b Board) MarshalJSON() ([]byte, error) {
	b.mut.RLock()
	defer b.mut.RUnlock()
	return json.Marshal(b.Posts)
}

// FilterBoard returns a new Board with a subset of posts filtered by a string
func (b Board) FilterBoard(queries []string) *Board {
	b.mut.RLock()
	defer b.mut.RUnlock()
	selectedPosts := []Post{}
	for _, post := range b.Posts {
		postData := strings.ToLower(post.Title)
		mismatch := false
		for _, query := range queries {
			if !strings.Contains(postData, query) {
				mismatch = true
				break
			}
		}
		if !mismatch {
			selectedPosts = append(selectedPosts, post)
		}
	}
	board := NewBoard(selectedPosts)
	return &board
}

// GetPostByID returns a post that matches the postID
func (b Board) GetPostByID(postID int64) *Post {
	b.mut.RLock()
	defer b.mut.RUnlock()
	for _, post := range b.Posts {
		if post.ID == postID {
			return &post
		}
	}
	return nil
}

// LimitBoard modifies the current board with maxResults posts starting at offset
func (b *Board) LimitBoard(offset, maxResults int) {
	b.mut.Lock()
	defer b.mut.Unlock()
	if offset > len(b.Posts) {
		offset = len(b.Posts)
	}
	endIndex := offset + maxResults
	if endIndex > len(b.Posts) {
		endIndex = len(b.Posts)
	}
	b.Posts = b.Posts[offset:endIndex]
}

// SortPostsByLikes sorts Posts in reverse number of likes order
func (b *Board) SortPostsByLikes() {
	b.mut.Lock()
	defer b.mut.Unlock()
	b.sortPostsByLikes()
}

func (b *Board) sortPostsByLikes() {
	sort.Sort(sort.Reverse(SortByLikes(b.Posts)))
}

// SortByLikes is an interface for Sorting
type SortByLikes []Post

func (a SortByLikes) Len() int           { return len(a) }
func (a SortByLikes) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a SortByLikes) Less(i, j int) bool { return a[i].Likes < a[j].Likes }

// RandomizePosts will shuffle the current Board's posts
func (b *Board) RandomizePosts() {
	b.mut.Lock()
	defer b.mut.Unlock()
	rand.Shuffle(len(b.Posts), func(i, j int) {
		b.Posts[i], b.Posts[j] = b.Posts[j], b.Posts[i]
	})
}

// URLs returns an array of URLs of all the posts
func (b Board) URLs() []string {
	b.mut.RLock()
	defer b.mut.RUnlock()
	urls := []string{}
	for _, post := range b.Posts {
		urls = append(urls, post.InternalURL())
	}
	return urls
}

// Keywords returns the most popular words in posts
func (b Board) Keywords() []string {
	b.mut.RLock()
	defer b.mut.RUnlock()
	words := map[string]int{}
	for _, post := range b.Posts {
		for _, word := range strings.Fields(post.Title) {
			if len(word) > 4 {
				words[word]++
			}
		}
	}
	// Reuse Board sort
	board := NewBoard([]Post{})
	for word, count := range words {
		board.AddPost(Post{Title: word, Likes: int64(count)})
	}
	board.SortPostsByLikes()
	count := 0
	keywords := []string{}
	for _, post := range board.Posts {
		keywords = append(keywords, post.Title)
		count++
		if count >= MaxKeywords {
			break
		}
	}
	return keywords
}
