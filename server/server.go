package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

//LaunchServerOn Listens for incoming http connections on given port
func LaunchServerOn(port string) {

	http.HandleFunc("/task", HandleTask)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

//HandleTask handles POST request, understood as tasks, with well-formed payload
func HandleTask(w http.ResponseWriter, r *http.Request) {

	reqTime := initialize()

	if r.Method != "POST" {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var es Task
	err = json.Unmarshal(reqBody, &es)

	if err != nil {
		http.Error(w, "Internal Server Error: could not unmarshal POST body", http.StatusInternalServerError)
		return
	}

	if es.Execute.URL == "" || es.Execute.Body == "" {
		http.Error(w, "Bad Reqeust", http.StatusBadRequest)
		return
	}

	requester := Requester{
		client: es.Execute,
	}

	resp, err := requester.Request()

	if err != nil {
		http.Error(w, "Could not make POST request to upstream service", http.StatusInternalServerError)
		return
	}

	event := generateEvent(reqTime, rand.Int(), resp)
	mEvent, err := marshalEvent(event)

	if err != nil {
		http.Error(w, "Internal Server Error: could not unmarshal POST body", http.StatusInternalServerError)
		return
	}

	w.Write(mEvent)

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
