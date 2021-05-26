package fiware

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestRoadSegment(t *testing.T) {
	location := [][2]float64{{1.0, 1.0}, {2.0, 2.0}}
	ts := time.Now()
	rs := NewRoadSegment("segid", "segname", "roadid", location, &ts)

	if rs.ID != "urn:ngsi-ld:RoadSegment:segid" {
		t.Error(fmt.Sprintf("Expectation failed. Road id %s != %s", rs.ID, "segid"))
	}

	probability := 0.8
	surfaceType := "snow"
	rs = rs.WithSurfaceType(surfaceType, probability)

	if rs.SurfaceType.Probability != probability {
		t.Error(fmt.Sprintf("Expectation failed. Surface type probability %f != %f", rs.SurfaceType.Probability, probability))
	}

	if rs.SurfaceType.Value != surfaceType {
		t.Error(fmt.Sprintf("Expectation failed. Surface type %s != %s", rs.SurfaceType.Value, surfaceType))
	}
}

func TestDeviceModel(t *testing.T) {
	categories := []string{"temperature"}

	deviceModel := NewDeviceModel("id", categories)

	if deviceModel.Category.Value[0] != categories[0] {
		t.Error(fmt.Sprintf("Expectation failed. Category does not match %s", categories))
	}

	// test devicemodel categories are as expected
}

func TestTrafficFlowObserved(t *testing.T) {
	id := TrafficFlowObservedIDPrefix + "trafficFlowObservedID"
	location := [2]float64{1.0, 1.0}
	ts := time.Now().String()

	tfo := NewTrafficFlowObserved(id, location[0], location[1], ts, 1)
	if tfo == nil {
		t.Error("Expectation failed. TrafficFlowObserved is empty")
	}

}
