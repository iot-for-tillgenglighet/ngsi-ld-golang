package ngsi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/geojson"
)

func TestRegisterContextSource(t *testing.T) {
	registrationBody, _ := NewCsourceRegistration("Point", []string{"x", "y"}, "lolcathost", nil)
	jsonBytes, _ := json.Marshal(registrationBody)
	ctxRegistry := NewContextRegistry()
	req, _ := http.NewRequest("POST", createURL("/csourceRegistration"), bytes.NewBuffer(jsonBytes))
	w := httptest.NewRecorder()

	NewRegisterContextSourceHandler(ctxRegistry).ServeHTTP(w, req)

	q, _ := newQueryFromParameters(req, []string{"Point"}, []string{"x"}, "")
	sources := ctxRegistry.GetContextSourcesForQuery(q)

	if len(sources) != 1 {
		t.Error("The registered context source was not added to the registry.")
	}

	if w.Code != http.StatusCreated {
		t.Error("Wrong status code returned. ", w.Code, " != expected 201")
	}
}

func TestRegisterContextSourceWithIDPatternMatch(t *testing.T) {
	regex := fmt.Sprintf("^%s.+", fiware.DeviceIDPrefix)
	registrationBody, _ := NewCsourceRegistration("A", []string{"a"}, "lolcathost", &regex)
	jsonBytes, _ := json.Marshal(registrationBody)
	ctxRegistry := NewContextRegistry()
	req, _ := http.NewRequest("POST", createURL("/csourceRegistration"), bytes.NewBuffer(jsonBytes))
	w := httptest.NewRecorder()

	NewRegisterContextSourceHandler(ctxRegistry).ServeHTTP(w, req)

	sources := ctxRegistry.GetContextSourcesForEntity(fiware.DeviceIDPrefix + "mydevice")

	if len(sources) != 1 {
		t.Error("The registered context source was not added to the registry.")
	}

	if w.Code != http.StatusCreated {
		t.Error("Wrong status code returned. ", w.Code, " != expected 201")
	}
}

func TestThatRequestsWithIDPatternMatchAreForwardedToRemoteContext(t *testing.T) {
	mockService := setupMockServiceThatReturns(204, "application/ld+json", "")
	defer mockService.Close()

	remoteURL := mockService.URL
	regex := "urn:ngsi-ld:TypeA:.+"
	registrationBody, _ := NewCsourceRegistration("TypeA", []string{"a"}, remoteURL, &regex)
	jsonBytes, _ := json.Marshal(registrationBody)
	ctxRegistry := NewContextRegistry()

	// Send a POST request to register a remote context source
	req, _ := http.NewRequest("POST", createURL("/csourceRegistration"), bytes.NewBuffer(jsonBytes))
	w := httptest.NewRecorder()
	NewRegisterContextSourceHandler(ctxRegistry).ServeHTTP(w, req)

	// Send a PATCH request to update entity attributes (that are handled by the "remote" source)
	entityID := "urn:ngsi-ld:TypeA:myentity"
	req, _ = http.NewRequest("PATCH", "https://localhost/ngsi-ld/v1/entities/"+entityID+"/attrs/", nil)
	request := newRequestWrapper(req)
	sources := ctxRegistry.GetContextSourcesForEntity(entityID)

	for _, src := range sources {
		err := src.UpdateEntityAttributes(entityID, request)
		if err != nil {
			t.Error("Failed with unexpected error", err.Error())
			return
		}
	}
}

func TestThatRequestsAreForwardedToRemoteContext(t *testing.T) {
	mockService := setupMockServiceThatReturns(200, "application/ld+json", snowHeightResponseJSON)
	defer mockService.Close()

	remoteURL := mockService.URL
	registrationBody, _ := NewCsourceRegistration("WeatherObserved", []string{"snowHeight"}, remoteURL, nil)
	jsonBytes, _ := json.Marshal(registrationBody)
	ctxRegistry := NewContextRegistry()

	// Send a POST request to register a remote context source
	req, _ := http.NewRequest("POST", createURL("/csourceRegistration"), bytes.NewBuffer(jsonBytes))
	w := httptest.NewRecorder()
	NewRegisterContextSourceHandler(ctxRegistry).ServeHTTP(w, req)

	// Send a GET request for entities of type WeatherObserved (that are handled by the "remote" source)
	req, _ = http.NewRequest("GET", "https://localhost/ngsi-ld/v1/entities?type=WeatherObserved", nil)
	query, _ := newQueryFromParameters(req, []string{"WeatherObserved"}, []string{"snowHeight"}, "")
	sources := ctxRegistry.GetContextSourcesForQuery(query)

	numEntities := 0

	for _, src := range sources {
		src.GetEntities(query, func(entity Entity) error {
			numEntities++
			return nil
		})
	}

	if numEntities == 0 {
		t.Error("Failed to get entities from remote endpoint.")
	}
}

func TestThatGeoJSONResponsesAreProperlyPropagated(t *testing.T) {
	mockService := setupMockServiceThatReturns(200, geojson.ContentType, beachResponseJSON)
	defer mockService.Close()

	remoteURL := mockService.URL
	registrationBody, _ := NewCsourceRegistration("Beach", []string{""}, remoteURL, nil)
	jsonBytes, _ := json.Marshal(registrationBody)
	ctxRegistry := NewContextRegistry()

	// Send a POST request to register a remote context source
	req, _ := http.NewRequest("POST", createURL("/csourceRegistration"), bytes.NewBuffer(jsonBytes))
	w := httptest.NewRecorder()
	NewRegisterContextSourceHandler(ctxRegistry).ServeHTTP(w, req)

	// Send a GET request for entities of type Beach (that are handled by the "remote" source)
	req, _ = http.NewRequest("GET", "https://localhost/ngsi-ld/v1/entities?type=Beach", nil)
	query, _ := newQueryFromParameters(req, []string{"Beach"}, []string{""}, "")
	sources := ctxRegistry.GetContextSourcesForQuery(query)

	numEntities := 0

	for _, src := range sources {
		src.GetEntities(query, func(entity Entity) error {
			numEntities++
			bytes, _ := json.Marshal(entity)
			fmt.Printf("entity: %s", string(bytes))
			return nil
		})
	}

	if numEntities == 0 {
		t.Error("Failed to get entities from remote endpoint.")
	}
}

const beachResponseJSON string = `{"type": "FeatureCollection","features": [
	{"id":"urn:ngsi-ld:Beach:42","type": "Feature",
	"geometry": {
		"type": "MultiPolygon",
		"coordinates": [[[
			[16.826877016818194,62.371366230256456],[16.82746858045308,62.37197792385098],
			[16.826075957396505,62.37229386059263],[16.825800236618605,62.37160561482045],
			[16.826877016818194,62.371366230256456]
			]]]
	},
	"properties": {
	  "description": "En fin liten strand.",
	  "location": {
		"type": "MultiPolygon",
		"coordinates": [[[
			  [16.826877016818194,62.371366230256456],[16.82746858045308,62.37197792385098],
			  [16.826075957396505,62.37229386059263],[16.825800236618605,62.37160561482045],
			  [16.826877016818194,62.371366230256456]
			]]]
	  },
	  "name": "Stranden",
	  "refSeeAlso": [
		"urn:ngsi-ld:Device:tempsensor-19"
	  ],
	  "type": "Beach"
	}}]}`

const snowHeightResponseJSON string = "[{\"id\": \"urn:ngsi-ld:WeatherObserved:SnowHeight:snow_10a52aaa84c35727:2020-04-08T15:01:32Z\", \"type\": \"WeatherObserved\",\"dateObserved\": { \"type\": \"Property\", \"value\": {\"@type\": \"DateTime\", \"@value\": \"2020-04-08T15:01:32Z\"}}, \"location\": { \"type\": \"GeoProperty\", \"value\": { \"type\": \"Point\", \"coordinates\": [16.5687632, 62.4081681]}}, \"refDevice\": {\"type\": \"Relationship\", \"object\": \"urn:ngsi-ld:Device:snow_10a52aaa84c35727\"}, \"snowHeight\": { \"type\": \"Property\", \"value\": 0}, \"@context\": [\"https://schema.lab.fiware.org/ld/context\", \"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld\"]}]"

func setupMockServiceThatReturns(responseCode int, contentType, body string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", contentType)
		w.WriteHeader(responseCode)
		if body != "" {
			w.Write([]byte(body))
		}
	}))
}
