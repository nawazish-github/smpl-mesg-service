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

func Test_ShouldTestClientsPOSTMethod(t *testing.T) {
	requester := Requester{
		client: StubExecute{},
	}

	b, c, e := requester.Request()

	if e != nil {
		t.Errorf("Requester failed to make POST call")
	}

	strResp := string(b)
	if strResp != "some data" {
		t.Errorf("Client could not return valid response")
	}

	if c != 200 {
		t.Errorf("Client could not return StatusOK")
	}

}

func (r StubExecute) POST() ([]byte, int, error) {
	b := []byte("some data")
	return b, 200, nil
}

func Test_ShouldReturnValidResponseAfterRunningBatchRequests(t *testing.T) {
	sExecute := []Execute{{URL: "url1", Body: "body1"}, {URL: "url2", Body: "body2"}}
	var batchExecute BatchExecute
	batchExecute.BatchRequests = sExecute
	br := StubBatchRequester{batchExecute: batchExecute, client: StubExecute{}}
	bites, _, _ := br.Request()
	var m map[string]Success
	err := json.Unmarshal(bites, &m)
	if err != nil {
		t.Errorf("Batch executions failed.")
	}
	if len(m) != 2 {
		t.Errorf("Batch executions did not execute exactly 2 POST requests")
	}
}

func (br StubBatchRequester) Request() ([]byte, int, error) {
	var m = make(map[string]Success)
	for _, v := range br.batchExecute.BatchRequests {
		_, statuscode, _ := br.client.POST()
		m[v.URL] = Success{
			StatusCode: statuscode,
		}
	}
	b, e := json.Marshal(m)
	return b, 200, e
}

func (r stubBatchExecute) POST() ([]byte, int, error) {
	b := []byte("some data")
	return b, 200, nil
}

type StubBatchRequester struct {
	batchExecute BatchExecute
	client       Client
}
