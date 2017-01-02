package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func getTemplateDir() string {
	return os.Getenv("SERVER_TEMPLATES")
}

var indexPath = fmt.Sprintf("%s/index.htm", getTemplateDir())
var uRLFilePaths = map[string]string{
	"/": indexPath,
}

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

func handler(w http.ResponseWriter, r *http.Request) {
	// Execute the template per HTTP request
	filePath, err := getFilePath(r.URL.Path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	data, err := readFile(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintf(w, data)
}

// Run starts up the HTTP server
func Run() {
	address := ":" + os.Getenv("PORT")
	fmt.Println("server listening on", address)
	http.HandleFunc("/", handler)
	http.ListenAndServe(address, nil)
}
