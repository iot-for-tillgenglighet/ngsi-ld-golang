package geojson

import (
	"encoding/json"
	"fmt"
	"reflect"
)

const (
	ContentType            string = "application/geo+json"
	ContentTypeWithCharset string = "application/geo+json;charset=utf-8"
)

type GeoJSONFeatureProperty struct {
}

type GeoJSONFeature interface {
	SetProperty(name string, value interface{})
}

type GeoJSONGeometry interface {
	GeoPropertyType() string
	GeoPropertyValue() GeoJSONGeometry
}

type geoJSONGeometryImpl struct {
	Geometry GeoJSONGeometry
}

func (gjgi *geoJSONGeometryImpl) MarshalJSON() ([]byte, error) {
	return json.Marshal(gjgi.Geometry)
}

func (gjgi *geoJSONGeometryImpl) UnmarshalJSON(data []byte) error {
	temp := struct {
		Type        string          `json:"type"`
		Coordinates json.RawMessage `json:"coordinates"`
	}{}

	err := json.Unmarshal(data, &temp)
	if err != nil {
		return err
	}

	if temp.Type == "MultiPolygon" {
		coords := [][][][]float64{}
		err = json.Unmarshal(temp.Coordinates, &coords)
		if err != nil {
			return err
		}

		gjp := CreateGeoJSONPropertyFromMultiPolygon(coords)
		gjgi.Geometry = gjp.Value
	} else if temp.Type == "Point" {
		coords := [2]float64{}
		err = json.Unmarshal(temp.Coordinates, &coords)
		if err != nil {
			return err
		}

		gjp := CreateGeoJSONPropertyFromWGS84(coords[0], coords[1])
		gjgi.Geometry = gjp.Value
	} else {
		return fmt.Errorf("unable to unmarshal geometry of type %s", temp.Type)
	}

	return nil
}

type geoJSONFeatureImpl struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Geometry   geoJSONGeometryImpl    `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type GeoJSONFeatureCollection struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
	Context  *[]string        `json:"@context,omitempty"`
}

func (gjfc *GeoJSONFeatureCollection) UnmarshalJSON(data []byte) error {
	collection := struct {
		Features []geoJSONFeatureImpl `json:"features"`
	}{}

	err := json.Unmarshal(data, &collection)
	if err != nil {
		return err
	}

	for idx := range collection.Features {
		gjfc.Features = append(gjfc.Features, &collection.Features[idx])
	}

	return nil
}

type SpatialEntity interface {
	ToGeoJSONFeature(propertyName string, simplified bool) (GeoJSONFeature, error)
}

func NewGeoJSONFeature(id, typ string, geom GeoJSONGeometry) GeoJSONFeature {
	f := &geoJSONFeatureImpl{
		ID:       id,
		Type:     "Feature",
		Geometry: geoJSONGeometryImpl{Geometry: geom},
		Properties: map[string]interface{}{
			"type": typ,
		},
	}
	return f
}

func NewGeoJSONFeatureCollection(features []GeoJSONFeature, includeContext bool) *GeoJSONFeatureCollection {
	collection := &GeoJSONFeatureCollection{
		Type:     "FeatureCollection",
		Features: features,
	}

	if includeContext {
		collection.Context = &[]string{
			"https://schema.lab.fiware.org/ld/context",
			"https://uri.etsi.org/ngsi-ld/v1/ngsi-ld-core-context.jsonld",
		}
	}

	return collection
}

// TODO: Explain the peculiarities of nil interfaces to Go newcomers ...
func propertyIsNotNil(v interface{}) bool {
	if v == nil {
		return false
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return !reflect.ValueOf(v).IsNil()
	}
	return true
}

func (f *geoJSONFeatureImpl) SetProperty(name string, value interface{}) {
	if propertyIsNotNil(value) {
		if f.Properties == nil {
			f.Properties = map[string]interface{}{}
		}
		f.Properties[name] = value
	}
}

func NewEntityConverter(property string, simplified bool, collection *GeoJSONFeatureCollection) func(interface{}) interface{} {
	return func(e interface{}) interface{} {
		switch v := e.(type) {
		// Do not double convert features when they come from a remote source
		case GeoJSONFeature:
			collection.Features = append(collection.Features, v)
			return e
		// Certain entity types support a conversion to a GeoJSON feature ...
		case SpatialEntity:
			f, _ := v.(SpatialEntity).ToGeoJSONFeature(property, simplified)
			collection.Features = append(collection.Features, f)
			return f
		// ... and some dont. How can we handle those in a better way?
		default:
			return &geoJSONFeatureImpl{Type: "Feature"}
		}
	}
}

func UnpackGeoJSONToCallback(bytes []byte, callback func(GeoJSONFeature) error) error {
	collection := GeoJSONFeatureCollection{}
	err := json.Unmarshal(bytes, &collection)
	if err != nil {
		return err
	}

	for _, f := range collection.Features {
		err = callback(f)
		if err != nil {
			return err
		}
	}

	return nil
}

type Property struct {
	Type string `json:"type"`
}

//GeoJSONPropertyPoint is used as the value object for a GeoJSONPropertyPoint
type GeoJSONPropertyPoint struct {
	Type        string     `json:"type"`
	Coordinates [2]float64 `json:"coordinates"`
}

func (gjpp *GeoJSONPropertyPoint) GeoPropertyType() string {
	return gjpp.Type
}

func (gjpp *GeoJSONPropertyPoint) GeoPropertyValue() GeoJSONGeometry {
	return gjpp
}

//GeoJSONPropertyMultiPolygon is used as the value object for a GeoJSONPropertyPoint
type GeoJSONPropertyMultiPolygon struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates"`
}

func (gjpmp *GeoJSONPropertyMultiPolygon) GeoPropertyType() string {
	return gjpmp.Type
}

func (gjpmp *GeoJSONPropertyMultiPolygon) GeoPropertyValue() GeoJSONGeometry {
	return gjpmp
}

//GeoJSONProperty is used to encapsulate different GeoJSONGeometry types
type GeoJSONProperty struct {
	Property
	Value GeoJSONGeometry `json:"value"`
}

func (gjp *GeoJSONProperty) GeoPropertyType() string {
	return gjp.Value.GeoPropertyType()
}

func (gjp *GeoJSONProperty) GeoPropertyValue() GeoJSONGeometry {
	return gjp.Value
}

func CreateGeoJSONPropertyFromJSON(data []byte) *GeoJSONProperty {
	tmp := struct {
		Property
		Value geoJSONGeometryImpl `json:"value"`
	}{}

	err := json.Unmarshal(data, &tmp)
	if err != nil {
		return nil
	}

	p := &GeoJSONProperty{
		Property: tmp.Property,
		Value:    tmp.Value.Geometry,
	}

	return p
}

//CreateGeoJSONPropertyFromWGS84 creates a GeoJSONProperty from a WGS84 coordinate
func CreateGeoJSONPropertyFromWGS84(longitude, latitude float64) *GeoJSONProperty {
	p := &GeoJSONProperty{
		Property: Property{Type: "GeoProperty"},
		Value: &GeoJSONPropertyPoint{
			Type:        "Point",
			Coordinates: [2]float64{longitude, latitude},
		},
	}

	return p
}

//CreateGeoJSONPropertyFromMultiPolygon creates a GeoJSONProperty from an array of polygon coordinate arrays
func CreateGeoJSONPropertyFromMultiPolygon(coordinates [][][][]float64) *GeoJSONProperty {
	p := &GeoJSONProperty{
		Property: Property{Type: "GeoProperty"},
		Value: &GeoJSONPropertyMultiPolygon{
			Type:        "MultiPolygon",
			Coordinates: coordinates,
		},
	}

	return p
}
