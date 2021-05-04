package ngsi

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/datamodels/fiware"
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/geojson"
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

func TestGetBeachAsGeoJSON(t *testing.T) {
	req, _ := http.NewRequest("GET", createURL("/entitites?type=Beach"), nil)
	req.Header["Accept"] = []string{"application/geo+json"}

	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("Beach", "temperature")

	location := geojson.CreateGeoJSONPropertyFromMultiPolygon([][][][]float64{})
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

func TestGetLongLatForGeoPropertyPoint(t *testing.T) {
	location := geojson.CreateGeoJSONPropertyFromWGS84(17.2961, 65.2789)

	beach := fiware.NewBeach("beach", "Beach Beach", location).WithDescription("A very beachy beach.")
	fmt.Print(beach.Location.GetAsPoint())
}

func TestGetLongLatForGeoPropertyMultiPolygon(t *testing.T) {
	poly := [][][][]float64{}
	row1 := [][][]float64{}
	row2 := [][]float64{}
	row3 := []float64{17.2961, 65.2789}
	row2 = append(row2, row3)
	row1 = append(row1, row2)
	poly = append(poly, row1)

	location := geojson.CreateGeoJSONPropertyFromMultiPolygon(poly)

	beach := fiware.NewBeach("beach", "Beach Beach", location).WithDescription("A very beachy beach.")
	fmt.Print(beach.Location.GetAsPoint())
}

func TestGetWaterQualityObservedAsGeoJSON(t *testing.T) {
	req, _ := http.NewRequest("GET", createURL("/entitites?type=WaterQualityObserved&options=keyValues"), nil)
	req.Header["Accept"] = []string{"application/geo+json"}

	w := httptest.NewRecorder()
	contextRegistry := NewContextRegistry()
	contextSource := newMockedContextSource("WaterQualityObserved", "temperature")

	wqo1 := fiware.NewWaterQualityObserved("badtempsensor", 64.2789, 17.2961, "2021-04-22T17:23:41Z")

	contextSource.entities = append(
		contextSource.entities,
		wqo1,
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

	location := geojson.CreateGeoJSONPropertyFromMultiPolygon([][][][]float64{})
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
