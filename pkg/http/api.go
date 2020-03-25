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
	"fmt"
	"strconv"

	"github.com/AndreasM009/eventstore-go/store"
	jsoniter "github.com/json-iterator/go"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

// APIRoutes to register routes
type APIRoutes interface {
	RegisterRoutes(r *routing.Router)
}

type api struct {
	evtstore store.EventStore
	json     jsoniter.API
}

const (
	entityIDParam     = "id"
	versionQueryParam = "version"
)

// NewAPI creates a new server instance
func NewAPI(evtstore store.EventStore) APIRoutes {
	api := &api{
		evtstore: evtstore,
		json:     jsoniter.ConfigFastest,
	}
	return api
}

func (a *api) RegisterRoutes(r *routing.Router) {
	r.Post("/entities/<id>", a.onPostEntity)
	r.Put("/entities/<id>", a.onPutEntity)
	r.Get("/entities/<id>", a.onGetEntity)
}

func (a *api) onPostEntity(c *routing.Context) error {
	id := c.Param(entityIDParam)
	body := c.PostBody()

	ety := store.Entity{}

	err := a.json.Unmarshal(body, &ety)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_POST_ENTITY", fmt.Sprintf("can't deserialize request: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	ety.ID = id
	ety.Version = 0

	res, err := a.evtstore.Add(&ety)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_POST_ENTITY", fmt.Sprintf("can't append entity to eventstore: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	resdata, err := a.json.Marshal(res)
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
	body := c.PostBody()

	ety := store.Entity{}

	err := a.json.Unmarshal(body, &ety)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_PUT_ENTITY", fmt.Sprintf("can't deserialize request: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	ety.ID = id

	res, err := a.evtstore.Append(&ety)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_PUT_ENTITY", fmt.Sprintf("can't add entity to eventstore: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	resdata, err := a.json.Marshal(res)
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
	vstr := c.QueryArgs().Peek(versionQueryParam)

	if vstr != nil {
		v, err := strconv.ParseInt(string(vstr), 10, 64)

		if err != nil {
			msg := NewErrorResponse("ERR_INVOKE_GET_ENTITY", fmt.Sprintf("can't convert version to number: %s", err))
			respondWithError(c.RequestCtx, fasthttp.StatusBadRequest, msg)
			return nil
		}

		version = v
	} else {
		v, err := a.evtstore.GetLatestVersionNumber(id)
		if err != nil {
			msg := NewErrorResponse("ERR_INVOKE_GET_ENTITY", fmt.Sprintf("can't get latest version: %s", err))
			respondWithError(c.RequestCtx, fasthttp.StatusNotFound, msg)
			return nil
		}

		version = v
	}

	ety, err := a.evtstore.GetByVersion(id, version)

	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_GET_ENTITY", fmt.Sprintf("can't load entity: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusNotFound, msg)
		return nil
	}

	resdata, err := a.json.Marshal(ety)
	if err != nil {
		msg := NewErrorResponse("ERR_INVOKE_GET_ENTITY", fmt.Sprintf("can't serialize to respond: %s", err))
		respondWithError(c.RequestCtx, fasthttp.StatusInternalServerError, msg)
		return nil
	}

	respondWithJSON(c.RequestCtx, fasthttp.StatusOK, resdata)
	return nil
}
