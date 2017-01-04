package server

import (
	"testing"
)

func TestGetURLHandler(t *testing.T) {
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
	handler, err = getURLHandler("/asdf")
	if err == nil {
		t.Fail()
	}
	if handler == nil {
		t.Fail()
	}
	data, err = handler()
	if err == nil {
		t.Fail()
	}
}
