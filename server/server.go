package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

//LaunchServerOn Listens for incoming http connections on given port
func LaunchServerOn(port string) {

	http.HandleFunc("/task", HandleTask)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

//HandleTask handles POST request, understood as tasks, with well-formed payload
func HandleTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var es task
	err = json.Unmarshal(reqBody, &es)

	if err != nil {
		http.Error(w, "Internal Server Error: could not unmarshal POST body", http.StatusInternalServerError)
		return
	}

	if es.Execute.URL == "" || es.Execute.Body == "" {
		http.Error(w, "Bad Reqeust", http.StatusBadRequest)
		return
	}

}
