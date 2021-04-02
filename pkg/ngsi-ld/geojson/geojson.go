package geojson

import "reflect"

type GeoJSONFeatureProperty struct {
}

type GeoJSONFeature interface {
	SetProperty(name string, value interface{})
}

type GeoJSONGeometry interface {
	GeoPropertyType() string
	GeoPropertyValue() GeoJSONGeometry
}

type geoJSONFeatureImpl struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Geometry   GeoJSONGeometry        `json:"geometry"`
	Properties map[string]interface{} `json:"properties"`
}

type GeoJSONFeatureCollection struct {
	Type     string           `json:"type"`
	Features []GeoJSONFeature `json:"features"`
	Context  *[]string        `json:"@context,omitempty"`
}

type SpatialEntity interface {
	ToGeoJSONFeature(propertyName string, simplified bool) (GeoJSONFeature, error)
}

func NewGeoJSONFeature(id, typ string, geom GeoJSONGeometry) GeoJSONFeature {
	f := &geoJSONFeatureImpl{
		ID:       id,
		Type:     "Feature",
		Geometry: geom,
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
		f.Properties[name] = value
	}
}

func NewEntityConverter(property string, simplified bool, collection *GeoJSONFeatureCollection) func(interface{}) interface{} {
	return func(e interface{}) interface{} {
		switch v := e.(type) {
		case SpatialEntity:
			f, _ := v.(SpatialEntity).ToGeoJSONFeature(property, simplified)
			collection.Features = append(collection.Features, f)
			return f
		default:
			return &geoJSONFeatureImpl{Type: "Feature"}
		}
	}
}
