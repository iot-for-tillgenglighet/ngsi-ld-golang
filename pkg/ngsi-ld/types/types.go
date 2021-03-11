package types

import (
	"strconv"
)

//BaseEntity contains the required base properties an Entity must have
type BaseEntity struct {
	ID      string   `json:"id"`
	Type    string   `json:"type"`
	Context []string `json:"@context"`
}

//Property contains the mandatory Type property
type Property struct {
	Type string `json:"type"`
}

//DateTimeProperty stores date and time values (surprise, surprise ...)
type DateTimeProperty struct {
	Property
	Value struct {
		Type  string `json:"@type"`
		Value string `json:"@value"`
	} `json:"value"`
}

//CreateDateTimeProperty creates a property from a UTC time stamp
func CreateDateTimeProperty(value string) *DateTimeProperty {
	dtp := &DateTimeProperty{
		Property: Property{Type: "Property"},
	}

	dtp.Value.Type = "DateTime"
	dtp.Value.Value = value

	return dtp
}

//GeoJSONProperty is used to store lat/lon coordinates
type GeoJSONProperty struct {
	Property
	Value struct {
		Type        string     `json:"type"`
		Coordinates [2]float64 `json:"coordinates"`
	} `json:"value"`
}

//CreateGeoJSONPropertyFromWGS84 creates a GeoJSONProperty from a WGS84 coordinate
func CreateGeoJSONPropertyFromWGS84(longitude, latitude float64) GeoJSONProperty {
	p := GeoJSONProperty{
		Property: Property{Type: "GeoProperty"},
	}

	p.Value.Type = "Point"
	p.Value.Coordinates[0] = longitude
	p.Value.Coordinates[1] = latitude

	return p
}

//NumberProperty holds a float64 Value
type NumberProperty struct {
	Property
	Value float64 `json:"value"`
}

//NewNumberProperty is a convenience function for creating NumberProperty instances
func NewNumberProperty(value float64) *NumberProperty {
	return &NumberProperty{
		Property: Property{Type: "Property"},
		Value:    value,
	}
}

//NewNumberPropertyFromInt accepts a value as an int and returns a new NumberProperty
func NewNumberPropertyFromInt(value int) *NumberProperty {
	return &NumberProperty{
		Property: Property{Type: "Property"},
		Value:    float64(value),
	}
}

//TextProperty stores values of type text
type TextProperty struct {
	Property
	Value string `json:"value"`
}

//TextListProperty stores values of type text list
type TextListProperty struct {
	Property
	Value []string `json:"value"`
}

//NewTextListProperty accepts a value as a string array and returns a new TextListProperty
func NewTextListProperty(value []string) *TextListProperty {
	return &TextListProperty{
		Property: Property{Type: "Property"},
		Value:    value,
	}
}

//NewNumberPropertyFromString accepts a value as a string and returns a new NumberProperty
func NewNumberPropertyFromString(value string) *NumberProperty {
	number, _ := strconv.ParseFloat(value, 64)
	return &NumberProperty{
		Property: Property{Type: "Property"},
		Value:    number,
	}
}

//NewTextProperty accepts a value as a string and returns a new TextProperty
func NewTextProperty(value string) *TextProperty {
	return &TextProperty{
		Property: Property{Type: "Property"},
		Value:    value,
	}
}

//Relationship is a base type for all types of relationships
type Relationship struct {
	Type string `json:"type"`
}

//SingleObjectRelationship stores information about an entity's relation to a single object
type SingleObjectRelationship struct {
	Relationship
	Object string `json:"object"`
}

//NewSingleObjectRelationship accepts an object ID as a string and returns a new SingleObjectRelationship
func NewSingleObjectRelationship(object string) *SingleObjectRelationship {
	return &SingleObjectRelationship{
		Relationship: Relationship{Type: "Relationship"},
		Object:       object,
	}
}

//MultiObjectRelationship stores information about an entity's relation to multiple objects
type MultiObjectRelationship struct {
	Relationship
	Object []string `json:"object"`
}

//NewMultiObjectRelationship accepts an array of object ID:s and returns a new MultiObjectRelationship
func NewMultiObjectRelationship(objects []string) MultiObjectRelationship {
	p := MultiObjectRelationship{
		Relationship: Relationship{Type: "Relationship"},
	}

	p.Object = objects

	return p
}

type RoadSegmentLocation struct {
	Property
	Value struct {
		Type        string       `json:"type"`
		Coordinates [][2]float64 `json:"coordinates"`
	} `json:"value"`
}

func NewRoadSegmentLocation(roadCoords [][2]float64) RoadSegmentLocation {
	r := RoadSegmentLocation{
		Property: Property{Type: "GeoProperty"},
	}

	r.Value.Type = "LineString"
	r.Value.Coordinates = roadCoords

	return r
}
