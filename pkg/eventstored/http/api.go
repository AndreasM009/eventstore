package http

//---------------------------------------------------------------------------------------------
// APIServer for event store
// Routes:
// POST /entities -> creates a new entity
// PUT /entities/{id} -> adds a new entity version
// GET /entities/{id}?version={versionnumber} -> gets an entity with specified version
// GET /entities/{id} -> gets the latest version available for specified entity
//---------------------------------------------------------------------------------------------

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/AndreasM009/eventstore-go/store"
	"github.com/AndreasM009/eventstore-service-go/pkg/eventstored/config"
	registry "github.com/AndreasM009/eventstore-service-go/pkg/eventstored/eventstore"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// APIRoutes to register routes
type APIRoutes interface {
	RegisterRoutes(r *routing.Router)
}

type api struct {
	evtstores map[string]store.EventStore
	registry  registry.Registry
}

const (
	eventstoreNameParam = "name"
	entityIDParam       = "id"
	versionQueryParam   = "version"
)

// NewAPI creates a new server instance
func NewAPI(evtstores map[string]store.EventStore, registry registry.Registry) APIRoutes {
	api := &api{
		evtstores: evtstores,
		registry:  registry,
	}
	return api
}

func (a *api) RegisterRoutes(r *routing.Router) {
	r.Post("/eventstores/<name>/entities/<id>", a.onPostEntity)
	r.Put("/eventstores/<name>/entities/<id>", a.onPutEntity)
	r.Get("/eventstores/<name>/entities/<id>", a.onGetEntity)
	r.Post("/configurations/<name>", a.onPostConfiguration)
}

func (a *api) onPostEntity(c *routing.Context) error {
	id := c.Param(entityIDParam)
	name := c.Param(eventstoreNameParam)
	body := c.PostBody()

	fmt.Println(string(body))

	eventstore, ok := a.evtstores[name]
	if !ok {
		msg := NewErrorResponse("ERR_INVOKE_POST_ENTITY", fmt.Sprintf("Evenstore %s not found", name))
		respondWithError(c.RequestCtx, fasthttp.StatusNotFound, msg)
		return nil
	}

	ety := store.Entity{}

	err := json.Unmarshal(body, &ety)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_POST_ENTITY", fmt.Sprintf("can't deserialize request: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	ety.ID = id
	ety.Version = 0

	res, err := eventstore.Add(&ety)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_POST_ENTITY", fmt.Sprintf("can't append entity to eventstore: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	resdata, err := json.Marshal(res)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_POST_ENTITY", fmt.Sprintf("can't serialize to respond: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	lurl := fmt.Sprintf("/entities/%s", res.ID)
	respondWithJSON(c.RequestCtx, fasthttp.StatusCreated, resdata)
	c.RequestCtx.Response.Header.Add("Location", lurl)
	return nil
}

func (a *api) onPutEntity(c *routing.Context) error {
	id := c.Param(entityIDParam)
	name := c.Param(eventstoreNameParam)
	body := c.PostBody()

	eventstore, ok := a.evtstores[name]
	if !ok {
		msg := NewErrorResponse("ERR_INVOKE_PUT_ENTITY", fmt.Sprintf("Evenstore %s not found", name))
		respondWithError(c.RequestCtx, fasthttp.StatusNotFound, msg)
		return nil
	}

	ety := store.Entity{}

	err := json.Unmarshal(body, &ety)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_PUT_ENTITY", fmt.Sprintf("can't deserialize request: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	ety.ID = id

	res, err := eventstore.Append(&ety)
	if err != nil {
		evterr, ok := err.(store.EventStoreError)

		if !ok {
			msg := NewErrorResponse("ERR_INVOKE_PUT_ENTITY", fmt.Sprintf("can't add entity to eventstore: %s", err))
			respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
			return nil
		}

		if evterr.ErrorType == store.EntityNotFound {
			msg := NewErrorResponse("ERR_INVOKE_PUT_ENTITY", fmt.Sprintf("can't add entity to eventstore: %s", err))
			respondWithError(c.RequestCtx, fasthttp.StatusNotFound, msg)
			return nil
		}

		if evterr.ErrorType == store.VersionConflict {
			msg := NewErrorResponse("ERR_INVOKE_PUT_ENTITY", fmt.Sprintf("can't add entity to eventstore: %s", err))
			respondWithError(c.RequestCtx, fasthttp.StatusConflict, msg)
			return nil
		}

		msg := NewErrorResponse("ERR_INVOKE_PUT_ENTITY", fmt.Sprintf("can't add entity to eventstore: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	resdata, err := json.Marshal(res)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_PUT_ENTITY", fmt.Sprintf("can't serialize to respond: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	respondWithJSON(c.RequestCtx, fasthttp.StatusOK, resdata)
	return nil
}

func (a *api) onGetEntity(c *routing.Context) error {
	var version int64

	id := c.Param(entityIDParam)
	name := c.Param(eventstoreNameParam)
	vstr := c.QueryArgs().Peek(versionQueryParam)

	eventstore, ok := a.evtstores[name]
	if !ok {
		msg := NewErrorResponse("ERR_INVOKE_GET_ENTITY", fmt.Sprintf("Evenstore %s not found", name))
		respondWithError(c.RequestCtx, fasthttp.StatusNotFound, msg)
		return nil
	}

	if vstr != nil {
		v, err := strconv.ParseInt(string(vstr), 10, 64)

		if err != nil {
			msg := NewErrorResponse("ERR_INVOKE_GET_ENTITY", fmt.Sprintf("can't convert version to number: %s", err))
			respondWithError(c.RequestCtx, fasthttp.StatusBadRequest, msg)
			return nil
		}

		version = v
	} else {
		v, err := eventstore.GetLatestVersionNumber(id)
		if err != nil {
			msg := NewErrorResponse("ERR_INVOKE_GET_ENTITY", fmt.Sprintf("can't get latest version: %s", err))
			respondWithError(c.RequestCtx, fasthttp.StatusNotFound, msg)
			return nil
		}

		version = v
	}

	ety, err := eventstore.GetByVersion(id, version)

	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_GET_ENTITY", fmt.Sprintf("can't load entity: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusNotFound, msg)
		return nil
	}

	resdata, err := json.Marshal(ety)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_GET_ENTITY", fmt.Sprintf("can't serialize to respond: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	respondWithJSON(c.RequestCtx, fasthttp.StatusOK, resdata)
	return nil
}

func (a *api) onPostConfiguration(c *routing.Context) error {
	name := c.Param(eventstoreNameParam)
	body := c.PostBody()

	_, ok := a.evtstores[name]

	if !ok {
		// not my configuration
		respondWithStatus(c.RequestCtx, fasthttp.StatusOK)
		return nil
	}

	cfg := config.Configuration{}
	err := json.Unmarshal(body, &cfg)

	if err != nil {
		respondWithStatus(c.RequestCtx, fasthttp.StatusInternalServerError)
		log.Printf("api: configuration can't be deserialized: %s", err)
		return nil
	}

	s, err := a.registry.Create(cfg)
	if err != nil {
		log.Printf("api: failed to update store from configuration: %s", err)
		respondWithStatus(c.RequestCtx, fasthttp.StatusInternalServerError)
		return nil
	}

	a.evtstores[cfg.Metadata.Name] = s
	log.Printf("api: configuration for Eventstore %s updated", cfg.Metadata.Name)
	respondWithStatus(c.RequestCtx, fasthttp.StatusOK)
	return nil
}
