package server

//IPOSTRequest interface is a generic interface to make all sorts of
//out bound calls from this server. An implementation of this
//interface - a "Requester" - in conjunction with a "Client",
// makes an appropriate outbound call.
type IPOSTRequest interface {
	Request() ([]byte, error)
}

//Request is an implementation of the IPOSTRequest contract
func (r Requester) Request() ([]byte, error) {
	b, e := r.client.POST()
	return b, e
}
