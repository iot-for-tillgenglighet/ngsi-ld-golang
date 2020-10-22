package ngsi

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/errors"
	"github.com/iot-for-tillgenglighet/ngsi-ld-golang/pkg/ngsi-ld/types"
)

//Entity is an informational representative of something that is supposed to exist in the real world, physically or conceptually
type Entity interface {
}

//QueryEntitiesCallback is used when queried context sources should pass back any
//entities matching the query that has been passed in
type QueryEntitiesCallback func(entity Entity) error

//NewQueryEntitiesHandler handles GET requests for NGSI entitites
func NewQueryEntitiesHandler(ctxReg ContextRegistry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		entityTypeNames := r.URL.Query().Get("type")
		attributeNames := r.URL.Query().Get("attrs")

		if entityTypeNames == "" && attributeNames == "" {
			errors.ReportNewBadRequestData(
				w,
				"A request for entities MUST specify at least one of type or attrs.",
			)
			return
		}

		entityTypes := strings.Split(entityTypeNames, ",")
		attributes := strings.Split(attributeNames, ",")

		q := r.URL.Query().Get("q")
		query := newQueryFromParameters(r, entityTypes, attributes, q)

		contextSources := ctxReg.GetContextSourcesForQuery(query)

		var entities = []Entity{}
		var err error

		for _, source := range contextSources {
			err = source.GetEntities(query, func(entity Entity) error {
				entities = append(entities, entity)
				return nil
			})
			if err != nil {
				break
			}
		}

		if err != nil {
			errors.ReportNewInternalError(
				w,
				"An internal error was encountered when trying to get entities from the context source.",
			)
			return
		}

		bytes, err := json.MarshalIndent(entities, "", "  ")
		if err != nil {
			errors.ReportNewInternalError(w, "Failed to encode response.")
			return
		}

		w.Header().Add("Content-Type", "application/ld+json")
		w.Write(bytes)
	})
}

//NewUpdateEntityAttributesHandler handles PATCH requests for NGSI entitity attributes
func NewUpdateEntityAttributesHandler(ctxReg ContextRegistry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: Replace this string manipulation with a callback that can use the http router's
		//		 functionality to extract URL params ...
		entitiesIdx := strings.Index(r.URL.Path, "/entities/")
		attrsIdx := strings.LastIndex(r.URL.Path, "/attrs/")

		if entitiesIdx == -1 || attrsIdx == -1 || attrsIdx < entitiesIdx {
			errors.ReportNewBadRequestData(
				w,
				"The supplied URL is invalid.",
			)
			return
		}

		entityID := r.URL.Path[entitiesIdx+10 : attrsIdx]

		patch := newPatchFromParameters(r)
		contextSources := ctxReg.GetContextSourcesForEntity(entityID)

		if len(contextSources) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		for _, source := range contextSources {
			err := source.UpdateEntityAttributes(entityID, patch)
			if err != nil {
				errors.ReportNewInvalidRequest(w, "Unable to update entity attributes: "+err.Error())
				return
			}
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

//NewCreateEntityHandler handles incoming POST requests for NGSI entities
func NewCreateEntityHandler(ctxReg ContextRegistry) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := newPostFromParameters(r)

		entity := &types.BaseEntity{}
		request.DecodeBodyInto(entity)

		contextSources := ctxReg.GetContextSourcesForEntityType(entity.Type)

		if len(contextSources) == 0 {
			errors.ReportNewInvalidRequest(w, "No context sources found matching the provided type")
			return
		}

		for _, source := range contextSources {
			err := source.CreateEntity(entity.Type, entity.ID, request)
			if err != nil {
				errors.ReportNewInvalidRequest(w, "Failed to create entity: "+err.Error())
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
	})
}
