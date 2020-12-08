package ngsi

import (
	"fmt"
	"testing"
)

func TestGeoCoordinatesParser(t *testing.T) {
	coords, err := parseGeometryCoordinates("[[2.4,2.1],[3.3,3.7]]")
	if err != nil {
		t.Error("Got error from coordinates parser", err)
	}

	if len(coords) != 4 {
		t.Error("Expected 2 coordinates, got", len(coords)/2)
	} else {
		if coords[0] != 2.4 || coords[1] != 2.1 {
			t.Error(fmt.Sprintf("First position should be (2.4,2.1) and not (%f,%f)", coords[0], coords[1]))
		}
	}
}
