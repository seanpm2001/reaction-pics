package server

import (
	"testing"
)

func TestGetFilePath(t *testing.T) {
	filePath, err := getFilePath("/")
	if err != nil {
		t.Fail()
	}
	if filePath[len(filePath)-9:] != "index.htm" {
		t.Fail()
	}
	filePath, err = getFilePath("/asdf")
	if err == nil {
		t.Fail()
	}
	if filePath != "" {
		t.Fail()
	}
}
