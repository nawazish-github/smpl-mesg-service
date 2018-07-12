package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HandleTaskShouldReturnStatusMethodNotAllowed(t *testing.T) {

	req, err := http.NewRequest("GET", "/task", nil)

	if err != nil {
		t.Errorf("Could not create request: %v", err.Error())
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HandleTask)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("HandleTask is receiving non-POST requests as well!")
	}
}
