package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
)

func HandleBatchExecute(w http.ResponseWriter, r *http.Request) {
	reqTime := initialize()

	//serve only HTTP POST requests
	if r.Method != "POST" {
		log.Fatal("Unsupported Method")
		http.Error(w, "Could not parse batch requests", http.StatusUnsupportedMediaType)
		return
	}

	//parse HTTP request for batch requests
	reqBytes, err := ioutil.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
		http.Error(w, "Could not parse batch requests", http.StatusBadRequest)
		return
	}

	event := generateEvent(reqTime, rand.Int(), reqBytes)

	//send event to mesg-core
	go sendAnEventToMesgCore(event)

	//execute the batch of POST requests
	//marshal request body to get batch requests
	var be BatchExecute
	err = json.Unmarshal(reqBytes, &be)
	if err != nil {
		log.Fatal(err)
		http.Error(w, "Could not marshal batch requests", http.StatusInternalServerError)
		return
	}

	batchRequester := BatchRequester{
		batchExecute: be,
	}

	mapBytes, statusCode, err := batchRequester.Request()

	if statusCode != 200 || err != nil {
		log.Fatal(err)
		http.Error(w, "Failed to process batch request", http.StatusInternalServerError)
		return
	}

	w.Write(mapBytes)

}
