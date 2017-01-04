package server

import (
	"testing"
)

func TestGetRootURLHandler(t *testing.T) {
	handler, err := getURLHandler("/")
	if err != nil {
		t.Fail()
	}
	data, err := handler()
	if err != nil {
		t.Fail()
	}
	if len(data) == 0 {
		t.Fail()
	}
}

func TestGetDataURLHandler(t *testing.T) {
	handler, err := getURLHandler("/data.json")
	if err != nil {
		t.Fail()
	}
	if handler == nil {
		t.Fail()
	}
	data, err := handler()
	if err != nil {
		t.Fail()
	}
	if data != "null" {
		t.Fail()
	}
}

func TestGetUnknownURLHandler(t *testing.T) {
	handler, err := getURLHandler("/asdf")
	if err == nil {
		t.Fail()
	}
	if handler == nil {
		t.Fail()
	}
	data, err := handler()
	if err == nil {
		t.Fail()
	}
	if data != "" {
		t.Fail()
	}
}
