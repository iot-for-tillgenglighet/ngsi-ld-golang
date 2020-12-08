package ngsi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

func createURL(path string, params ...string) string {
	url := "http://localhost:8080/ngsi-ld/v1" + path

	if len(params) > 0 {
		url = url + "?"

		for _, p := range params {
			url = url + p + "&"
		}

		url = strings.TrimSuffix(url, "&")
	}

	return url
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestNewRoadSegment(t *testing.T) {
	id := "urn:ngsi-ld:RoadSegment:road1"
	roadID := "road"
	name := "segName"

	coords := [][2]float64{
		{0.0, 0.0},
		{1.1, 6.5},
	}

	roadSegment := fiware.NewRoadSegment(id, name, roadID, coords)

	if len(roadSegment.Location.Value.Coordinates) != 2 {
		t.Error("Number of coords not as expected.")
	}

	for index, point := range roadSegment.Location.Value.Coordinates {
		if point != coords[index] {
			t.Error("Coords do not match.")
		}
	}
}

func TestCreateEntityUsesCorrectTypeAndID(t *testing.T) {
	entityID := "urn:ngsi-ld:Device:livboj"
	byteReader, typeName := newEntityAsByteBuffer(entityID)
	req, _ := http.NewRequest("POST", createURL("/entities"), byteReader)
	w := httptest.NewRecorder()

	ctxReg, ctxSrc := newContextRegistryWithSourceForType(typeName)

	NewCreateEntityHandler(ctxReg).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Error("Handler did not return the expected status code. ", w.Code, " != ", http.StatusCreated)
	}

	if ctxSrc.createdEntityType != typeName {
		t.Error("CreateEntity called with wrong type name. ", ctxSrc.createdEntityType, " != ", typeName)
	}

	if ctxSrc.createdEntity != entityID {
		t.Error("CreateEntity called with wrong entity ID. ", ctxSrc.createdEntity, " != ", entityID)
	}
}

func TestCreateEntityFailsWithNoContextSources(t *testing.T) {
	byteBuffer, _ := newEntityAsByteBuffer("id")
	req, _ := http.NewRequest("POST", createURL("/entities"), byteBuffer)
	w := httptest.NewRecorder()

	NewCreateEntityHandler(NewContextRegistry()).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("Wrong response code when posting device with no context sources. ", w.Code, " is not ", http.StatusBadRequest)
	}
}

func TestCreateEntityHandlesFailureFromContextSource(t *testing.T) {
	byteBuffer, typeName := newEntityAsByteBuffer("id")
	req, _ := http.NewRequest("POST", createURL("/entities"), byteBuffer)
	w := httptest.NewRecorder()

	ctxReg, ctxSrc := newContextRegistryWithSourceForType(typeName)
	ctxSrc.createEntityShouldFailWithError = errors.New("failure")

	NewCreateEntityHandler(ctxReg).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("Wrong response code when posting device with no context sources. ", w.Code, " is not ", http.StatusBadRequest)
	}
}

func TestGetEntitiesWithoutAttributesOrTypesFails(t *testing.T) {
	req, _ := http.NewRequest("GET", createURL("/entitites"), nil)
	w := httptest.NewRecorder()

	NewQueryEntitiesHandler(NewContextRegistry()).ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Error("GET /entities MUST require either type or attrs request parameter")
	}
}

func TestGetEntitiesWithAttribute(t *testing.T) {
	req, _ := http.NewRequest("GET", createURL("/entitites", "attrs=snowHeight"), nil)
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextRegistry.Register(newMockedContextSource(
		"", "snowHeight",
		e(""),
	))

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Error("That did not work .... :(")
	}
}

func TestGetEntitiesForDevice(t *testing.T) {
	deviceID := "urn:ngsi-ld:Device:mydevice"
	req, _ := http.NewRequest("GET", createURL("/entitites", "attrs=snowHeight", "q=refDevice==\""+deviceID+"\""), nil)
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("", "snowHeight")
	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	if contextSource.queriedDevice != deviceID {
		t.Error("Queried device did not match expectations. ", contextSource.queriedDevice, " != ", deviceID)
	} else if w.Code != http.StatusOK {
		t.Error("That did not work ... :(")
	}
}

func TestGetEntitiesWithGeoQueryNearPoint(t *testing.T) {
	req, _ := http.NewRequest("GET", createURL(
		"/entitites",
		"type=RoadSegment",
		"georel=near;maxDistance==2000",
		"geometry=Point",
		"coordinates=[8,40]"),
		nil)
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("RoadSegment", "")
	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("Unexpected response code", w.Code, w.Body.String())
		return
	}

	query := contextSource.generatedQuery
	if query.IsGeoQuery() == false {
		t.Error("Expected a GeoQuery from the QueryEntititesHandler")
	} else {
		geo := query.Geo()

		if geo.GeoRel != "near" {
			t.Error("Geospatial relation not correctly saved in geo query (" + geo.GeoRel + " != near)")
		}

		distance, _ := geo.Distance()
		if distance != 2000 {
			t.Error("Unexpected near distance parsed from geo query:", distance, "!=", 2000)
		}

		x, y, _ := geo.Point()
		if x != 8 || y != 40 {
			t.Error("Mismatching point: (", x, ",", y, ") != ( 8 , 40 )")
		}
	}
}

func TestGetEntitiesWithGeoQueryWithinRect(t *testing.T) {
	req, _ := http.NewRequest("GET", createURL(
		"/entitites",
		"type=RoadSegment",
		"georel=within",
		"geometry=Polygon",
		"coordinates=[[8,40],[9,41],[10,42]]"),
		nil)
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("RoadSegment", "")
	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("Handler failed with exit code", w.Code, w.Body.String())
		return
	}

	query := contextSource.generatedQuery
	if query.IsGeoQuery() == false {
		t.Error("Expected a GeoQuery from the QueryEntititesHandler")
	} else {
		geo := query.Geo()

		if geo.GeoRel != "within" {
			t.Error("Geospatial relation not correctly saved in geo query (" + geo.GeoRel + " != within)")
		}

		lon0, lat0, lon1, lat1, _ := geo.Rectangle()
		if lon0 != 8 || lat0 != 40 || lon1 != 10 || lat1 != 42 {
			t.Error("Bad coordinates in GeoQuery rect")
		}
	}
}

func TestUpdateEntitityAttributes(t *testing.T) {
	deviceID := "urn:ngsi-ld:Device:mydevice"
	jsonBytes, _ := json.Marshal(e("testvalue"))

	req, _ := http.NewRequest("PATCH", createURL("/entities/"+deviceID+"/attrs/"), bytes.NewBuffer(jsonBytes))
	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("", "value")
	contextRegistry.Register(contextSource)

	NewUpdateEntityAttributesHandler(contextRegistry).ServeHTTP(w, req)

	if contextSource.patchedEntity != deviceID {
		t.Error("Patched entity did not match expectations. ", contextSource.patchedEntity, " != ", deviceID)
	}
}

type mockEntity struct {
	Value string
}

func e(val string) mockEntity {
	return mockEntity{Value: val}
}

func newContextRegistryWithSourceForType(typeName string) (ContextRegistry, *mockCtxSource) {
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource(typeName, "")
	contextRegistry.Register(contextSource)
	return contextRegistry, contextSource
}

func newEntityAsByteBuffer(entityID string) (io.Reader, string) {
	device := fiware.NewDevice(entityID, "")
	jsonBytes, _ := json.Marshal(device)
	return bytes.NewBuffer(jsonBytes), device.Type
}

func newMockedContextSource(typeName string, attributeName string, e ...mockEntity) *mockCtxSource {
	source := &mockCtxSource{typeName: typeName, attributeName: attributeName}
	for _, entity := range e {
		source.entities = append(source.entities, entity)
	}
	return source
}

type mockCtxSource struct {
	typeName      string
	attributeName string
	entities      []Entity

	createEntityShouldFailWithError error

	queriedDevice     string
	createdEntity     string
	createdEntityType string
	patchedEntity     string

	generatedQuery Query
}

func (s *mockCtxSource) CreateEntity(typeName, entityID string, r Request) error {
	if s.createEntityShouldFailWithError == nil {
		s.createdEntity = entityID
		s.createdEntityType = typeName

		entity := &types.BaseEntity{}
		return r.DecodeBodyInto(entity)
	}

	return s.createEntityShouldFailWithError
}

func (s *mockCtxSource) GetEntities(q Query, cb QueryEntitiesCallback) error {

	s.generatedQuery = q

	if q.HasDeviceReference() {
		s.queriedDevice = q.Device()
	}

	for _, e := range s.entities {
		cb(e)
	}
	return nil
}

func (s *mockCtxSource) UpdateEntityAttributes(entityID string, req Request) error {
	s.patchedEntity = entityID
	e := &mockEntity{}
	return req.DecodeBodyInto(e)
}

func (s *mockCtxSource) ProvidesAttribute(attributeName string) bool {
	return s.attributeName == attributeName
}

func (s *mockCtxSource) ProvidesEntitiesWithMatchingID(entityID string) bool {
	return true
}

func (s *mockCtxSource) ProvidesType(typeName string) bool {
	return s.typeName == typeName
}
