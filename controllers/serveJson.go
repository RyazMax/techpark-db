package controllers

import (
	"github.com/mailru/easyjson"
	"github.com/valyala/fasthttp"
)

func serveJson(ctx *fasthttp.RequestCtx, code int, data easyjson.Marshaler) {
	if data != nil {
		body, _ := easyjson.Marshal(data)
		ctx.SetContentType("application/json")
		ctx.Write(body)
	}
	ctx.SetStatusCode(code)
}
