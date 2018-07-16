package server

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Client interface {
	POST() ([]byte, int, error)
}

func (e Execute) POST() ([]byte, int, error) {
	resp, err := http.Post(e.URL, "application/json", strings.NewReader(e.Body))
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return b, resp.StatusCode, err
}
