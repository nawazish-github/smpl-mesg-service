package server

import (
	"log"
	"net/http"
)

func LaunchServerOn(port string) {

	http.HandleFunc("/task", HandleTask)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func HandleTask(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		http.Error(w, "Unsupported Method", http.StatusMethodNotAllowed)
		return
	}

}
