package ngsi

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
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
	req := p.Request()

	// Request bodies can only be read once, so read the request's body ...
	buf, _ := ioutil.ReadAll(req.Body)
	// ... and replace it with a new reader with the same contents ...
	req.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
	// ... before returning yet another new reader to the caller
	return ioutil.NopCloser(bytes.NewBuffer(buf))
}

func (p *patchWrapper) DecodeBodyInto(v interface{}) error {
	return json.NewDecoder(p.BodyReader()).Decode(v)
}
