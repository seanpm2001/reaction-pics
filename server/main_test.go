package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
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
	if response.Body.String() != "null" {
		t.Fail()
	}
}
