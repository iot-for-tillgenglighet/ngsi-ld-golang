package ngsi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

func TestGetEntitiesAsGeoJSON(t *testing.T) {
	req, _ := http.NewRequest("GET", createURL("/entitites?type=Beach"), nil)
	req.Header["Accept"] = []string{"application/geo+json"}

	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("Beach", "temperature")

	location := types.CreateGeoJSONPropertyFromMultiPolygon([][][][]float64{})
	b1 := fiware.NewBeach("omaha", "Omaha Beach", location).WithDescription("This is a nice beach!")
	b1.WaterTemperature = types.NewNumberProperty(7.2)

	contextSource.entities = append(
		contextSource.entities,
		b1,
	)

	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	fmt.Printf("Got response: %s\n", w.Body)

	if w.Code != 200 {
		t.Error("Failed to get geojson data")
	}
}

func TestGetEntitiesAsSimplifiedGeoJSON(t *testing.T) {
	req, _ := http.NewRequest("GET", createURL("/entitites?type=Beach&options=keyValues"), nil)
	req.Header["Accept"] = []string{"application/geo+json"}

	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("Beach", "temperature")

	location := types.CreateGeoJSONPropertyFromMultiPolygon([][][][]float64{})
	b1 := fiware.NewBeach("omaha", "Omaha Beach", location).WithDescription("This is a nice beach!")
	b1.WaterTemperature = types.NewNumberProperty(7.2)

	contextSource.entities = append(
		contextSource.entities,
		b1,
	)

	contextRegistry.Register(contextSource)

	NewQueryEntitiesHandler(contextRegistry).ServeHTTP(w, req)

	fmt.Printf("Got response: %s\n", w.Body)

	if w.Code != 200 {
		t.Error("Failed to get geojson data")
	}
}
