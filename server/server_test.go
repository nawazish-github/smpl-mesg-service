package server

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

func Test_HandleTaskShouldReceiveValidHTTPPOSTBody(t *testing.T) {
	payload1 := `{
					"execute":{
						"url":"",
						"body":""
					}
				}`
	payload2 := `{
					"execute":{
						"url":"",
						"body":"payload2"
					}
				}`
	payload3 := `{
					"execute":{
						"url":"http://someurl:8080/task/payload3",
						"body":""
					}
				}`

	payloads := []string{payload1, payload2, payload3}

	for _, p := range payloads {
		req, _ := http.NewRequest("POST", "/task", strings.NewReader(p))
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(HandleTask)

		handler.ServeHTTP(rr, req)

		bs, _ := ioutil.ReadAll(req.Body)

		var es Task
		json.Unmarshal(bs, &es)
		if rr.Code != http.StatusBadRequest {
			t.Errorf("Invalid POST Body to HandlerTask. URL: %v and Body: %v", es.Execute.URL, es.Execute.Body)
		}
	}
}
