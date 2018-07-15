package server

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//LaunchServerOn Listens for incoming http connections on given port
func LaunchServerOn(port string) {

	http.HandleFunc("/execute", HandleTask)
	http.HandleFunc("/batchexecute", HandleBatchExecute)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

//HandleTask handles POST request, understood as tasks, with well-formed payload
//These are single tasks.
func HandleTask(w http.ResponseWriter, r *http.Request) {

	reqTime := initialize()

	if r.Method != "POST" {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error: Could not read request",
			http.StatusInternalServerError)
		return
	}

	event := generateEvent(reqTime, rand.Int(), reqBody)
	go sendAnEventToMesgCore(event)

	var es Task
	err = json.Unmarshal(reqBody, &es)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error: could not unmarshal POST body",
			http.StatusInternalServerError)
		return
	}

	//check if request is well formed
	//URL and Body are mandatory
	if es.Execute.URL == "" || es.Execute.Body == "" {
		log.Fatal(err)
		http.Error(w, "Bad Reqeust", http.StatusBadRequest)
		return
	}

	requester := Requester{
		client: es.Execute,
	}

	resp, code, err := requester.Request()

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Could not make POST request to upstream service",
			http.StatusInternalServerError)
		return
	}

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error: could not unmarshal POST body",
			http.StatusInternalServerError)
		return
	}

	successMsg := generateSuccessMessage(resp, code)
	ssm, err := serializeSuccessMessage(successMsg)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Internal Server Error: Could not generate success message",
			http.StatusInternalServerError)
		return
	}

	w.Write(ssm)
}

func generateEvent(date time.Time, id int, body []byte) *Event {
	return &Event{
		OnRequest: OnRequest{
			Time: date,
			Id:   id,
			Body: body,
		},
	}
}

func marshalEvent(e *Event) ([]byte, error) {
	return json.Marshal(e)
}

func initialize() time.Time {
	reqTime := time.Now()
	rand.Seed(time.Now().UnixNano())
	return reqTime
}

func sendAnEventToMesgCore(event *Event) {
	b, err := json.Marshal(event)
	if err != nil {
		log.Fatal(err)
		return
	}
	resp, err := http.Post("http://mesgcore", "application/json", bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println(resp.StatusCode)
}

func generateSuccessMessage(body []byte, statusCode int) *Success {
	return &Success{
		StatusCode: statusCode,
		Body:       body,
	}
}

func serializeSuccessMessage(s *Success) ([]byte, error) {
	b, e := json.Marshal(s)
	return b, e
}
