package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/albertyw/reaction-pics/tumblr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type HandlerTestSuite struct {
	suite.Suite
	deps handlerDeps
}

func (s *HandlerTestSuite) SetupTest() {
	logger := zap.NewNop().Sugar()
	board := tumblr.NewBoard([]tumblr.Post{})
	s.deps = handlerDeps{
		logger:         logger,
		board:          &board,
		appCacheString: appCacheString(logger),
	}
}

func (s *HandlerTestSuite) TestIndexFile(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	indexHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 200)

	assert.Contains(t, response.Body.String(), s.deps.appCacheString)
}

func (s *HandlerTestSuite) TestOnlyIndexFile(t *testing.T) {
	request, err := http.NewRequest("GET", "/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	indexHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 404)
}

func (s *HandlerTestSuite) TestReadFile(t *testing.T) {
	request, err := http.NewRequest("GET", "/static/favicon/manifest.json", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	staticHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 200)
	assert.True(t, len(response.Body.String()) > 100)
}

func (s *HandlerTestSuite) TestNoExactURL(t *testing.T) {
	request, err := http.NewRequest("GET", "/static/asdf.js", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	staticHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 404)

	response = httptest.NewRecorder()
	indexHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 404)
}

func (s *HandlerTestSuite) TestSearchHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/search", nil)
	assert.NoError(t, err)

	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"data\":[],\"offset\":0,\"totalResults\":0}")
}

func (s *HandlerTestSuite) TestSearchHandlerOffset(t *testing.T) {
	request, err := http.NewRequest("GET", "/search?offset=1", nil)
	assert.NoError(t, err)

	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"data\":[],\"offset\":1,\"totalResults\":0}")
}

func (s *HandlerTestSuite) TestSearchHandlerMalformedOffset(t *testing.T) {
	request, err := http.NewRequest("GET", "/search?offset=asdf", nil)
	assert.NoError(t, err)

	q := request.URL.Query()
	q.Add("query", "searchTerm")
	response := httptest.NewRecorder()
	searchHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"data\":[],\"offset\":0,\"totalResults\":0}")
}

func (s *HandlerTestSuite) TestPostHandlerMalformed(t *testing.T) {
	request, err := http.NewRequest("GET", "/post/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 404)
}

func (s *HandlerTestSuite) TestPostHandlerNotFound(t *testing.T) {
	request, err := http.NewRequest("GET", "/post/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 404)
}

func (s *HandlerTestSuite) TestPostHandler(t *testing.T) {
	post := tumblr.Post{ID: 1234}
	s.deps.board.AddPost(post)
	request, err := http.NewRequest("GET", "/post/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 200)
	assert.NotEqual(t, len(response.Body.String()), 0)
}

func (s *HandlerTestSuite) TestPostDataHandler(t *testing.T) {
	post := tumblr.Post{ID: 1234}
	s.deps.board.AddPost(post)
	request, err := http.NewRequest("GET", "/postdata/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 200)
	assert.NotEqual(t, len(response.Body.String()), 0)
}

func (s *HandlerTestSuite) TestPostDataPercentHandler(t *testing.T) {
	post := tumblr.Post{ID: 1234, Title: `asdf% qwer`}
	s.deps.board.AddPost(post)
	request, err := http.NewRequest("GET", "/postdata/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request, s.deps)
	var data map[string][]map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &data)
	title := data["data"][0]["title"].(string)
	assert.Equal(t, `asdf% qwer`, title)
}

func (s *HandlerTestSuite) TestPostDataHandlerMalformed(t *testing.T) {
	request, err := http.NewRequest("GET", "/postdata/asdf", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 404)
}

func (s *HandlerTestSuite) TestPostDataHandlerUnknown(t *testing.T) {
	request, err := http.NewRequest("GET", "/postdata/1234", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	postDataHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 404)
}

func (s *HandlerTestSuite) TestStatsHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/stats.json", nil)
	assert.NoError(t, err)

	response := httptest.NewRecorder()
	statsHandler(response, request, s.deps)
	assert.Equal(t, response.Code, 200)
	assert.Equal(t, response.Body.String(), "{\"keywords\":[],\"postCount\":\"0\"}")
}
