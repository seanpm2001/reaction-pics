package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albertyw/reaction-pics/tumblr"
)

func TestReadFile(t *testing.T) {
	handler := rewriteFS(http.FileServer(http.Dir(staticPath)).ServeHTTP)
	request, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fail()
	}
	response := httptest.NewRecorder()
	handler(response, request)
	if response.Code != 200 {
		t.Fail()
	}
	if len(response.Body.String()) < 100 {
		t.Fail()
	}
}

func TestExactURL(t *testing.T) {
	handler := rewriteFS(http.FileServer(http.Dir(staticPath)).ServeHTTP)
	request, err := http.NewRequest("GET", "/asdf", nil)
	if err != nil {
		t.Fail()
	}
	response := httptest.NewRecorder()
	handler(response, request)
	if response.Code != 404 {
		t.Fail()
	}
}

func TestNoExactURL(t *testing.T) {
	handler := rewriteFS(http.FileServer(http.Dir(staticPath)).ServeHTTP)
	request, err := http.NewRequest("GET", "/asdf", nil)
	if err != nil {
		t.Fail()
	}
	response := httptest.NewRecorder()
	handler(response, request)
	if response.Code != 404 {
		t.Fail()
	}
}

func TestDataURLHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/data.json", nil)
	if err != nil {
		t.Fail()
	}
	response := httptest.NewRecorder()
	dataURLHandler(response, request)
	if response.Code != 200 {
		t.Fail()
	}
	if response.Body.String() != "[]" {
		t.Fail()
	}
}

func TestSearchHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/search", nil)
	if err != nil {
		t.Fail()
	}
	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request)
	if response.Code != 200 {
		t.Fail()
	}
	if response.Body.String() != "[]" {
		t.Fail()
	}
}

func TestPostHandlerMalformed(t *testing.T) {
	request, err := http.NewRequest("GET", "/post/asdf", nil)
	if err != nil {
		t.Fail()
	}
	response := httptest.NewRecorder()
	postHandler(response, request)
	if response.Code != 404 {
		t.Fail()
	}
}

func TestPostHandlerNotFound(t *testing.T) {
	request, err := http.NewRequest("GET", "/post/1234", nil)
	if err != nil {
		t.Fail()
	}
	response := httptest.NewRecorder()
	postHandler(response, request)
	if response.Code != 404 {
		t.Fail()
	}
}

func TestPostHandler(t *testing.T) {
	post := tumblr.Post{ID: 1234}
	board.AddPost(post)
	defer func() { board.Reset() }()
	request, err := http.NewRequest("GET", "/post/1234", nil)
	if err != nil {
		t.Fail()
	}
	response := httptest.NewRecorder()
	postHandler(response, request)
	if response.Code != 200 {
		t.Fail()
	}
	if len(response.Body.String()) == 0 {
		t.Fail()
	}
}

func TestStatsHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/stats.json", nil)
	if err != nil {
		t.Fail()
	}
	response := httptest.NewRecorder()
	statsHandler(response, request)
	if response.Code != 200 {
		t.Fail()
	}
	if response.Body.String() != "{\"postCount\":\"0\"}" {
		t.Fail()
	}
}

func TestAddPost(t *testing.T) {
	oldLength := len(board.Posts)
	post1 := tumblr.Post{ID: 1}
	post2 := tumblr.Post{ID: 2}
	post3 := tumblr.Post{ID: 3}
	postChan := make(chan tumblr.Post, 1)
	defer func() { close(postChan) }()
	go loadPosts(postChan)
	postChan <- post1
	postChan <- post2
	postChan <- post3
	if len(board.Posts) == oldLength {
		t.Fail()
	}
	if board.Posts[0].ID > board.Posts[1].ID {
		t.Fail()
	}
}
