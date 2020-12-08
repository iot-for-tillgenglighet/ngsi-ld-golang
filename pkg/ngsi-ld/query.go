package ngsi

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

//Query is an interface to be used when passing queries to context registries and sources
type Query interface {
	HasDeviceReference() bool
	Device() string

	IsGeoQuery() bool
	Geo() *GeoQuery

	EntityAttributes() []string
	EntityTypes() []string

	Request() *http.Request
}

const (
	//GeoSpatialRelationNearPoint describes a relation as a max or min distance from a Point
	GeoSpatialRelationNearPoint = "near"
	//GeoSpatialRelationWithinRect describes a relation as an overlapping polygon
	GeoSpatialRelationWithinRect = "within"
)

//GeoQuery contains information about a geo-query that may be used for subscriptions
//or when querying entitites
type GeoQuery struct {
	Geometry    string    `json:"geometry"`
	Coordinates []float64 `json:"coordinates"`
	GeoRel      string    `json:"georel"`
	GeoProperty *string   `json:"geoproperty,omitempty"`

	distance uint32
}

//Distance returns the required distance in meters from a near Point and a boolean
//flag indicating if it is inclusive or exclusive
func (gq *GeoQuery) Distance() (uint32, bool) {
	return gq.distance, true
}

//Point extracts the position in the enclosed geometry
func (gq *GeoQuery) Point() (float64, float64, error) {
	if len(gq.Coordinates) == 2 {
		return gq.Coordinates[0], gq.Coordinates[1], nil
	}

	return 0, 0, errors.New("invalid number of coordinates in GeoQuery for a Point geometry")
}

//Rectangle extracts the two positions (opposing corners) within the enclosed geometry
func (gq *GeoQuery) Rectangle() (float64, float64, float64, float64, error) {
	if len(gq.Coordinates) == 6 {
		//TODO: Use all coordinates and allow for more elaborate polygons. For now we just
		// 		pick the first and the third coordinates and form a rect from them
		return gq.Coordinates[0], gq.Coordinates[1], gq.Coordinates[4], gq.Coordinates[5], nil
	}

	return 0, 0, 0, 0, fmt.Errorf("invalid number of coordinates in GeoQuery for a Polygon geometry")
}

func newQueryFromParameters(req *http.Request, types []string, attributes []string, q string) (Query, error) {

	var err error

	const refDevicePrefix string = "refDevice==\""

	qw := &queryWrapper{request: req, types: types, attributes: attributes}

	if strings.HasPrefix(q, refDevicePrefix) {
		splitElems := strings.Split(q, "\"")
		qw.device = &splitElems[1]
	}

	georel := req.URL.Query().Get("georel")
	if len(georel) > 0 {
		qw.geoQuery, err = newGeoQueryFromHTTPRequest(georel, req)
	}

	return qw, err
}

func newGeoQueryFromHTTPRequest(georel string, req *http.Request) (*GeoQuery, error) {

	var err error

	if georel == GeoSpatialRelationNearPoint {
		geoQuery := &GeoQuery{Geometry: "Point", GeoRel: GeoSpatialRelationNearPoint}

		if req.URL.Query().Get("geometry") != "Point" {
			return nil, errors.New("The geospatial relationship near is only defined for the geometry type Point")
		}

		distanceString := req.URL.Query().Get("maxDistance")
		if len(distanceString) < 2 || strings.HasPrefix(distanceString, "=") == false {
			return nil, errors.New("Required parameter maxDistance missing or invalid")
		}

		distanceString = distanceString[1:]
		distance, err := strconv.Atoi(distanceString)

		if err != nil {
			return nil, errors.New("Failed to parse distance: " + err.Error())
		}

		if distance < 0 {
			return nil, errors.New("Distance value must be non negative")
		}

		geoQuery.distance = uint32(distance)

		geoQuery.Coordinates, err = parseGeometryCoordinates(req.URL.Query().Get("coordinates"))
		if err != nil {
			return nil, err
		}

		if len(geoQuery.Coordinates) != 2 {
			fmt.Printf("expected one position for a Point geometry, but got %d\n", len(geoQuery.Coordinates)/2)
			return nil, fmt.Errorf("expected one position for a Point geometry, but got %d", len(geoQuery.Coordinates)/2)
		}

		return geoQuery, nil
	} else if georel == GeoSpatialRelationWithinRect {
		geoQuery := &GeoQuery{Geometry: "Polygon", GeoRel: GeoSpatialRelationWithinRect}

		if req.URL.Query().Get("geometry") != "Polygon" {
			return nil, errors.New("The geospatial relationship \"within\" is only defined for the geometry type Polygon")
		}

		geoQuery.Coordinates, err = parseGeometryCoordinates(req.URL.Query().Get("coordinates"))
		if err != nil {
			return nil, err
		}

		if len(geoQuery.Coordinates) != 6 {
			return nil, fmt.Errorf("The geospatial relationship \"within\" is only implemented for the Polygon type with five positions describing a bounding rect, but %d positions were received", len(geoQuery.Coordinates)/2)
		}

		return geoQuery, nil
	}

	return nil, errors.New("Only the geo-spatial relationships \"near\" and \"within\" are supported at this time")
}

func parseGeometryCoordinates(coordparameter string) ([]float64, error) {
	coordinates := []float64{}

	const init = 0
	const prelon = 1
	const lonint = 2
	const londec = 3
	const prelat = 4
	const latint = 5
	const latdec = 6

	state := init

	lon := 0.0
	lat := 0.0

	intpart := 0
	decpart := 0.0
	decfactor := 1.0

	pdepth := 0

	for i := range coordparameter {
		b := coordparameter[i]

		//fmt.Printf("Byte: %c State: %d\n", b, state)

		if state == init {
			if b != '[' {
				return nil, errors.New("coordinates string must start with a [")
			}
			state = prelon
			pdepth++
			continue
		}

		if b == '[' {
			if state != prelon {
				return nil, fmt.Errorf("Unexpected [ at position %d in %s", i, coordparameter)
			}
			pdepth++
		} else if b == ']' {
			pdepth--

			if pdepth < 0 {
				return nil, fmt.Errorf("Unexpected ] at position %d in %s", i, coordparameter)
			}

			if state == latint || state == latdec {
				lat = float64(intpart) + decpart
				intpart = 0
				decpart = 0.0

				coordinates = append(coordinates, lon, lat)
				state = prelon
			}

			if pdepth == 0 {
				break
			}
		} else if b >= '0' && b <= '9' {
			if state == prelon {
				state = lonint
			} else if state == prelat {
				state = latint
			}

			if state == latint || state == lonint {
				intpart = (intpart * 10) + int(b-'0')
			} else if state == latdec || state == londec {
				decfactor /= 10.0
				decpart += float64(b-'0') * decfactor
			}
		} else if b == '.' {
			if state == latint {
				state = latdec
			} else if state == lonint {
				state = londec
			}
		} else if b == ',' {
			if state == lonint || state == londec {
				lon = float64(intpart) + decpart
				intpart = 0
				decpart = 0.0
				state = prelat
			}
		} else {
			return nil, fmt.Errorf("Invalid byte '%c' found at position %d in %s", b, i, coordparameter)
		}
	}

	if pdepth > 0 {
		return nil, fmt.Errorf("Missing ] at end of coordinates array %s", coordparameter)
	}

	return coordinates, nil
}

type queryWrapper struct {
	request    *http.Request
	types      []string
	attributes []string
	device     *string

	geoQuery *GeoQuery
}

func (q *queryWrapper) HasDeviceReference() bool {
	return q.device != nil
}

func (q *queryWrapper) IsGeoQuery() bool {
	return q.Geo() != nil
}

func (q *queryWrapper) Geo() *GeoQuery {
	return q.geoQuery
}

func (q *queryWrapper) Device() string {
	return *q.device
}

func (q *queryWrapper) EntityAttributes() []string {
	return q.attributes
}

func (q *queryWrapper) EntityTypes() []string {
	return q.types
}

func (q *queryWrapper) Request() *http.Request {
	return q.request
}
