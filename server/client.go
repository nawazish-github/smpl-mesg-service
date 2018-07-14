package server

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type Client interface {
	POST() ([]byte, error)
}

func (e Execute) POST() ([]byte, error) {
	resp, err := http.Post(e.URL, "application/json", strings.NewReader(e.Body))
	b, err := ioutil.ReadAll(resp.Body)
	return b, err
}
