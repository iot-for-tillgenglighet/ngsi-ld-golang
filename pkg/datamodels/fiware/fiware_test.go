package fiware

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestRoadSegment(t *testing.T) {
	location := [][2]float64{{1.0, 1.0}, {2.0, 2.0}}
	rs := NewRoadSegment("segid", "segname", "roadid", location)

	if rs.ID != "urn:ngsi-ld:RoadSegment:segid" {
		t.Error(fmt.Sprintf("Expectation failed. Road id %s != %s", rs.ID, "segid"))
	}

	probability := 80.0
	surfaceType := "snow"
	rs = rs.WithSurfaceType(surfaceType, probability)

	if rs.SurfaceType.Probability != probability {
		t.Error(fmt.Sprintf("Expectation failed. Surface type probability %f != %f", rs.SurfaceType.Probability, probability))
	}

	if rs.SurfaceType.Value != surfaceType {
		t.Error(fmt.Sprintf("Expectation failed. Surface type %s != %s", rs.SurfaceType.Value, surfaceType))
	}
}
