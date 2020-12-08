package ngsi

import (
	"testing"
)

func TestGeoCoordinatesParser(t *testing.T) {
	coords, err := parseGeometryCoordinates("[[2.4,2],[3,3.7]]")
	if err != nil {
		t.Error("Got error from coordinates parser", err)
	}

	if len(coords) != 4 {
		t.Error("Expected 2 coordinates, got", len(coords)/2)
	} else {
		if coords[0] != 2.4 {
			t.Error("First longitude should be", 2.4, "not", coords[0])
		}
	}
}
