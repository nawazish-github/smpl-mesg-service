package server

type task struct {
	Execute execute `json:execute`
}

type execute struct {
	URL  string `json:url`
	Body string `json:body`
}
