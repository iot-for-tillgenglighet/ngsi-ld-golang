package ngsi

import (
	"bytes"
	"encoding/json"
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

func TestCreateEntityUsesCorrectTypeAndID(t *testing.T) {
	entityID := "urn:ngsi-ld:Device:livboj"
	device := fiware.NewDevice(entityID, "")
	jsonBytes, _ := json.Marshal(device)
	req, _ := http.NewRequest("POST", createURL("/entities"), bytes.NewBuffer(jsonBytes))
	w := httptest.NewRecorder()

	contextRegistry := NewContextRegistry()
	typeName := "Device"
	contextSource := newMockedContextSource(typeName, "")
	contextRegistry.Register(contextSource)

	NewCreateEntityHandler(contextRegistry).ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Error("Handler did not return the expected status code. ", w.Code, " != ", http.StatusCreated)
	}

	if contextSource.createdEntityType != typeName {
		t.Error("CreateEntity called with wrong type name. ", contextSource.createdEntityType, " != ", typeName)
	}

	if contextSource.createdEntity != entityID {
		t.Error("CreateEntity called with wrong entity ID. ", contextSource.createdEntity, " != ", entityID)
	}
}

func TestCreateEntityFailsWithNoContextSources(t *testing.T) {
	device := fiware.NewDevice("urn:ngsi-ld:Device:livboj", "on")
	jsonBytes, _ := json.Marshal(device)
	req, _ := http.NewRequest("POST", createURL("/entities"), bytes.NewBuffer(jsonBytes))
	w := httptest.NewRecorder()

	NewCreateEntityHandler(NewContextRegistry()).ServeHTTP(w, req)

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

	queriedDevice     string
	createdEntity     string
	createdEntityType string
	patchedEntity     string
}

func (s *mockCtxSource) CreateEntity(typeName, entityID string, post Post) error {
	s.createdEntity = entityID
	s.createdEntityType = typeName

	entity := &types.BaseEntity{}
	return post.DecodeBodyInto(entity)
}

func (s *mockCtxSource) GetEntities(q Query, cb QueryEntitiesCallback) error {

	if q.HasDeviceReference() {
		s.queriedDevice = q.Device()
	}

	for _, e := range s.entities {
		cb(e)
	}
	return nil
}

func (s *mockCtxSource) UpdateEntityAttributes(entityID string, patch Patch) error {
	s.patchedEntity = entityID
	e := &mockEntity{}
	return patch.DecodeBodyInto(e)
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
