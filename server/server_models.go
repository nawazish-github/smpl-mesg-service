package server

import "time"

type Task struct {
	Execute Execute `json:execute`
}

type Execute struct {
	URL  string `json:url`
	Body string `json:body`
}

type Event struct {
	OnRequest OnRequest `json:onRequest`
}

type OnRequest struct {
	Time time.Time `json:time`
	Id   int       `json:id`
	Body []byte    `json:body`
}

//StubExecute helps Unit Test Client.POST method
type StubExecute struct{}

//Requester holds an instance of Client implementation
//to make POST calls
type Requester struct {
	client Client
}

//Success is HTTP response body for the downstream callers
type Success struct {
	StatusCode int    `json:statusCode`
	Body       []byte `json:body`
}
