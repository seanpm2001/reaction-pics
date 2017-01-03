package server

import (
	"errors"
	"fmt"
	"github.com/albertyw/devops-reactions-index/tumblr"
	"io/ioutil"
	"net/http"
	"os"
)

func getTemplateDir() string {
	return os.Getenv("SERVER_TEMPLATES")
}

const dataURLPath = "/data.json"

var indexPath = fmt.Sprintf("%s/index.htm", getTemplateDir())
var uRLFilePaths = map[string]string{
	"/":         indexPath,
	dataURLPath: "",
}
var posts []tumblr.Post

func readFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	html := string(data)
	return html, nil
}

func getFilePath(urlPath string) (string, error) {
	filePath, exists := uRLFilePaths[urlPath]
	if !exists {
		return "", errors.New("")
	}
	return filePath, nil
}

func processURL(urlPath string) (string, error) {
	if urlPath == dataURLPath {
		return tumblr.PostsToJSON(posts), nil
	}
	return "", errors.New("no handler for url path")
}

func handler(w http.ResponseWriter, r *http.Request) {
	// Execute the template per HTTP request
	urlPath := r.URL.Path
	filePath, err := getFilePath(urlPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	var data string
	if filePath == "" {
		data, err = processURL(urlPath)
	} else {
		data, err = readFile(filePath)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, data)
}

// Run starts up the HTTP server
func Run(p []tumblr.Post) {
	posts = p
	address := ":" + os.Getenv("PORT")
	fmt.Println("server listening on", address)
	http.HandleFunc("/", handler)
	http.ListenAndServe(address, nil)
}
