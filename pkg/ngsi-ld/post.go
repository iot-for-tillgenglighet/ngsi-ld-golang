package ngsi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

//Post is an interface to be used when passing header and body information for Post requests
type Post interface {
	BodyReader() io.Reader
	DecodeBodyInto(v interface{}) error
	Request() *http.Request
}

func newPostFromParameters(req *http.Request) Post {
	pw := &postWrapper{request: req}
	return pw
}

type postWrapper struct {
	request *http.Request
}

func (p *postWrapper) Request() *http.Request {
	return p.request
}

func (p *postWrapper) BodyReader() io.Reader {
	req := p.Request()

	// Request bodies can only be read once, so read the request's body ...
	buf, _ := ioutil.ReadAll(req.Body)
	// ... and replace it with a new reader with the same contents ...
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	// ... before returning yet another new reader to the caller
	return ioutil.NopCloser(bytes.NewBuffer(buf))
}

func (p *postWrapper) DecodeBodyInto(v interface{}) error {
	return json.NewDecoder(p.BodyReader()).Decode(v)
}
