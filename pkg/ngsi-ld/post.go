package ngsi

import (
	"encoding/json"
	"io"
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
	return p.Request().Body
}

func (p *postWrapper) DecodeBodyInto(v interface{}) error {
	return json.NewDecoder(p.BodyReader()).Decode(v)
}
