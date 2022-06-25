package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	handler := NoSurf(&myH)

	// checking if handler is type of http.Handler
	switch v := handler.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Error(fmt.Printf("type is not http.Handler but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	handler := SessionLoad(&myH)
	// checking if handler is type of http.Handler
	switch v := handler.(type) {
	case http.Handler:
		//do nothing
	default:
		t.Error(fmt.Printf("type is not http.Handler but is %T", v))
	}
}
