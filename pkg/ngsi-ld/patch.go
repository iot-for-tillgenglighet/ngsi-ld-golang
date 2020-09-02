package ngsi

import (
	"encoding/json"
	"io"
	"net/http"
)

//Patch is an interface to be used when passing header and body information for PATCH requests
type Patch interface {
	BodyReader() io.Reader
	DecodeBodyInto(v interface{}) error
	Request() *http.Request
}

func newPatchFromParameters(req *http.Request) Patch {
	pw := &patchWrapper{request: req}
	return pw
}

type patchWrapper struct {
	request *http.Request
}

func (p *patchWrapper) Request() *http.Request {
	return p.request
}

func (p *patchWrapper) BodyReader() io.Reader {
	return p.Request().Body
}

func (p *patchWrapper) DecodeBodyInto(v interface{}) error {
	return json.NewDecoder(p.BodyReader()).Decode(v)
}
