package ngsi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

//Request is an interface to be used when passing header and body information for NGSI-LD API requests
type Request interface {
	BodyReader() io.Reader
	DecodeBodyInto(v interface{}) error
	Request() *http.Request
}

func newRequestWrapper(req *http.Request) Request {
	rw := &requestWrapper{request: req}
	return rw
}

type requestWrapper struct {
	request *http.Request
}

func (r *requestWrapper) Request() *http.Request {
	return r.request
}

func (r *requestWrapper) BodyReader() io.Reader {
	req := r.Request()

	// Request bodies can only be read once, so read the request's body ...
	buf, _ := ioutil.ReadAll(req.Body)
	// ... and replace it with a new reader with the same contents ...
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	// ... before returning yet another new reader to the caller
	return ioutil.NopCloser(bytes.NewBuffer(buf))
}

func (r *requestWrapper) DecodeBodyInto(v interface{}) error {
	return json.NewDecoder(r.BodyReader()).Decode(v)
}
