package http

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

const (
	jsonContentTypeHeader = "application/json"
)

func respond(ctx *fasthttp.RequestCtx, code int, obj []byte) {
	ctx.Response.SetStatusCode(code)
	ctx.Response.SetBody(obj)

	if len(ctx.Response.Header.ContentType()) == 0 {
		ctx.Response.Header.SetContentType(jsonContentTypeHeader)
	}
}

func respondWithJSON(ctx *fasthttp.RequestCtx, code int, obj []byte) {
	respond(ctx, code, obj)
	ctx.Response.Header.SetContentType(jsonContentTypeHeader)
}

func respondWithError(ctx *fasthttp.RequestCtx, code int, resp ErrorResponse) {
	b, _ := json.Marshal(&resp)
	respondWithJSON(ctx, code, b)
}
