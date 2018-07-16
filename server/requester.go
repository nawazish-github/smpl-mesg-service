package server

import (
	"encoding/json"
	"sync"
)

//IPOSTRequest interface is a generic interface to make all sorts of
//out bound calls from this server. An implementation of this
//interface - a "Requester" - in conjunction with a "Client",
// makes an appropriate outbound call.
type IPOSTRequest interface {
	Request() ([]byte, int, error)
}

//Request is an implementation of the IPOSTRequest contract
func (r Requester) Request() ([]byte, int, error) {
	b, c, e := r.client.POST()
	return b, c, e
}

//Request is an implementation of the IPOSTRequest contract.
//Implemented by BatchRequester. The fires concurrent POST
//requests using the client.
func (br *BatchRequester) Request() ([]byte, int, error) {
	successMap := make(map[string]Success)
	var wg sync.WaitGroup
	wg.Add(len(br.batchExecute.BatchRequests))

	for _, execute := range br.batchExecute.BatchRequests {
		go func() {
			defer wg.Done()
			respBytes, statusCode, _ := execute.POST()
			successMap[execute.URL] = Success{
				Body:       respBytes,
				StatusCode: statusCode,
			}
		}()
	}

	wg.Wait()
	mapBytes, err := json.Marshal(successMap)
	if err != nil {
		return nil, 500, err
	}
	return mapBytes, 200, nil
}
