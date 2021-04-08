package ngsi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"regexp"
	"strings"

	"github.com/google/uuid"
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/errors"
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/geojson"
)

//CsourceRegistration is a wrapper for information about a registered context source
type CsourceRegistration interface {
	Endpoint() string
	ProvidesAttribute(attributeName string) bool
	ProvidesEntitiesWithMatchingID(entityID string) bool
	ProvidesType(typeName string) bool
}

//NewRegisterContextSourceHandler handles POST requests for csource registrations
func NewRegisterContextSourceHandler(ctxReg ContextRegistry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := ioutil.ReadAll(r.Body)
		reg, err := NewCsourceRegistrationFromJSON(body)

		if err != nil {
			errors.ReportNewBadRequestData(
				w,
				"Failed to create registration from payload: "+err.Error(),
			)
			return
		}

		remoteCtxSrc, _ := NewRemoteContextSource(reg)

		ctxReg.Register(remoteCtxSrc)

		jsonBytes, _ := json.Marshal(remoteCtxSrc)

		w.WriteHeader(http.StatusCreated)
		w.Header().Add("Content-Type", "application/json")
		w.Write(jsonBytes)
	})
}

type remoteResponse struct {
	responseCode int
	headers      http.Header
	bytes        []byte
}

func (rr *remoteResponse) Header() http.Header {
	if rr.headers == nil {
		rr.headers = make(http.Header)
	}
	return rr.headers
}

func (rr *remoteResponse) MatchesContentType(contentType string) bool {
	responseType := rr.Header()["Content-Type"][0]
	return strings.HasPrefix(responseType, contentType)
}

func (rr *remoteResponse) Write(b []byte) (int, error) {
	rr.bytes = append(rr.bytes, b...)
	return len(b), nil
}

func (rr *remoteResponse) WriteHeader(responseCode int) {
	rr.responseCode = responseCode
}

//NewRemoteContextSource creates an instance of a ContextSource by wrapping a CsourceRegistration
func NewRemoteContextSource(registration CsourceRegistration) (ContextSource, error) {
	return &remoteContextSource{ID: uuid.New().String(), registration: registration}, nil
}

type remoteContextSource struct {
	ID           string `json:"id"`
	registration CsourceRegistration
}

func (rcs *remoteContextSource) CreateEntity(typeName, entityID string, r Request) error {
	u, _ := url.Parse(rcs.registration.Endpoint())
	req := r.Request()

	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme

	forwardedHost := req.Header.Get("Host")
	if forwardedHost != "" {
		req.Header.Set("X-Forwarded-Host", forwardedHost)
	}
	req.Host = u.Host

	// Change the User-Agent header to something more appropriate
	req.Header.Add("User-Agent", "ngsi-context-broker/0.1")

	response, err := proxyToRemote(u, req)

	if err != nil {
		return fmt.Errorf("attempt to create %s entity failed with status code %d: %s", typeName, response.responseCode, err.Error())
	}

	return err
}

func (rcs *remoteContextSource) GetEntities(query Query, callback QueryEntitiesCallback) error {
	u, _ := url.Parse(rcs.registration.Endpoint())
	req := query.Request()

	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme

	forwardedHost := req.Header.Get("Host")
	if forwardedHost != "" {
		req.Header.Set("X-Forwarded-Host", forwardedHost)
	}
	req.Host = u.Host

	// Change the User-Agent header to something more appropriate
	req.Header.Add("User-Agent", "ngsi-context-broker/0.1")

	// We do not want to propagate the Accept-Encoding header to prevent compression
	req.Header.Del("Accept-Encoding")

	response, err := proxyToRemote(u, req)

	// If the response code is 200 we can just unmarshal the payload
	// and pass the individual entitites to the supplied callback.
	// We need to check of the payload is GeoJSON or not though.
	if response.responseCode == http.StatusOK {

		if response.MatchesContentType(geojson.ContentType) {
			err = geojson.UnpackGeoJSONToCallback(response.bytes, func(f geojson.GeoJSONFeature) error {
				return callback(f)
			})
		} else {
			var unmarshaledResponse []interface{}
			err = json.Unmarshal(response.bytes, &unmarshaledResponse)
			if err == nil {
				for _, e := range unmarshaledResponse {
					callback(e)
				}
			}
		}
	}

	return err
}

func (rcs *remoteContextSource) UpdateEntityAttributes(entityID string, r Request) error {
	u, _ := url.Parse(rcs.registration.Endpoint())
	req := r.Request()

	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme

	forwardedHost := req.Header.Get("Host")
	if forwardedHost != "" {
		req.Header.Set("X-Forwarded-Host", forwardedHost)
	}
	req.Host = u.Host

	// Change the User-Agent header to something more appropriate
	req.Header.Add("User-Agent", "ngsi-context-broker/0.1")

	_, err := proxyToRemote(u, req)

	if err != nil {
		return fmt.Errorf("failed to patch entity %s: %s", entityID, err.Error())
	}

	return nil
}

func (rcs *remoteContextSource) ProvidesAttribute(attributeName string) bool {
	return rcs.registration.ProvidesAttribute(attributeName)
}

func (rcs *remoteContextSource) ProvidesEntitiesWithMatchingID(entityID string) bool {
	return rcs.registration.ProvidesEntitiesWithMatchingID(entityID)
}

func (rcs *remoteContextSource) ProvidesType(typeName string) bool {
	return rcs.registration.ProvidesType(typeName)
}

func (rcs *remoteContextSource) RetrieveEntity(entityID string, r Request) (Entity, error) {
	u, _ := url.Parse(rcs.registration.Endpoint())
	req := r.Request()

	req.URL.Host = u.Host
	req.URL.Scheme = u.Scheme

	forwardedHost := req.Header.Get("Host")
	if forwardedHost != "" {
		req.Header.Set("X-Forwarded-Host", forwardedHost)
	}
	req.Host = u.Host

	// Change the User-Agent header to something more appropriate
	req.Header.Add("User-Agent", "ngsi-context-broker/0.1")
	// We do not want to propagate the Accept-Encoding header to prevent compression
	req.Header.Del("Accept-Encoding")

	response, err := proxyToRemote(u, req)

	if err != nil {
		return nil, fmt.Errorf("failed to retrieve entity %s: %s", entityID, err.Error())
	}

	var entity interface{}

	if response.responseCode == http.StatusOK {
		err = json.Unmarshal(response.bytes, &entity)
		if err == nil {
			return entity, nil
		}

		return nil, fmt.Errorf(
			"failed to unmarshal retrieved entity %s from %s: %s",
			entityID, string(response.bytes), err.Error(),
		)
	}

	return nil, fmt.Errorf("unexpected response code from retrieve entity %s: %d != 200", entityID, response.responseCode)
}

func proxyToRemote(u *url.URL, req *http.Request) (remoteResponse, error) {
	response := remoteResponse{}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(&response, req)

	var err error

	if response.responseCode >= http.StatusBadRequest {
		if len(response.bytes) > 0 {
			err = fmt.Errorf("%s", string(response.bytes))
		} else {
			err = fmt.Errorf("received %d response with empty body", response.responseCode)
		}
	}

	return response, err
}

type ctxSrcReg struct {
	Type        string          `json:"type"`
	Information []ctxSrcRegInfo `json:"information"`
	Endpt       string          `json:"endpoint"`
}

func (csr *ctxSrcReg) Endpoint() string {
	return csr.Endpt
}

func (csr *ctxSrcReg) ProvidesAttribute(attributeName string) bool {
	for _, reginfo := range csr.Information {
		for _, attr := range reginfo.Properties {
			if attr == attributeName {
				return true
			}
		}
	}
	return false
}

func (csr *ctxSrcReg) ProvidesEntitiesWithMatchingID(entityID string) bool {
	for _, reginfo := range csr.Information {
		for _, entity := range reginfo.Entities {
			if entity.regexpForID != nil && entity.regexpForID.MatchString(entityID) {
				return true
			}
		}
	}
	return false
}

func (csr *ctxSrcReg) ProvidesType(typeName string) bool {
	for _, reginfo := range csr.Information {
		for _, entity := range reginfo.Entities {
			if entity.Type == typeName {
				return true
			}
		}
	}
	return false
}

type entityInfo struct {
	IDPattern   *string `json:"idPattern,omitempty"`
	regexpForID *regexp.Regexp

	Type string `json:"type"`
}

type ctxSrcRegInfo struct {
	Entities   []entityInfo `json:"entities"`
	Properties []string     `json:"properties"`
}

//NewCsourceRegistration creates and returns a concrete implementation of the CsourceRegistration interface
func NewCsourceRegistration(entityTypeName string, attributeNames []string, endpoint string, idpattern *string) (CsourceRegistration, error) {
	regInfo := ctxSrcRegInfo{Entities: []entityInfo{}, Properties: attributeNames}
	einfo := &entityInfo{Type: entityTypeName, IDPattern: idpattern}
	if idpattern != nil {
		var err error
		einfo.regexpForID, err = regexp.CompilePOSIX(*idpattern)
		if err != nil {
			return nil, err
		}
	}
	regInfo.Entities = append(regInfo.Entities, *einfo)

	reg := &ctxSrcReg{Type: "ContextSourceRegistration", Endpt: endpoint}
	reg.Information = []ctxSrcRegInfo{regInfo}

	return reg, nil
}

//NewCsourceRegistrationFromJSON unpacks a byte buffer into a CsourceRegistration and validates its contents
func NewCsourceRegistrationFromJSON(jsonBytes []byte) (CsourceRegistration, error) {
	registration := &ctxSrcReg{}
	err := json.Unmarshal(jsonBytes, registration)

	if err != nil {
		return nil, err
	}

	for infoIdx := range registration.Information {
		info := &registration.Information[infoIdx]

		for entityIdx := range info.Entities {
			entity := &info.Entities[entityIdx]

			if entity.IDPattern != nil {
				entity.regexpForID, err = regexp.CompilePOSIX(*entity.IDPattern)
				if err != nil {
					return nil, err
				}
			}
		}
	}

	// TODO: More validation ...

	return registration, nil
}
