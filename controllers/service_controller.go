package controllers

import (
	"net/http"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

func ServiceStatus(ctx *fasthttp.RequestCtx) {
	info := models.DatabaseInfo{}
	info.Get(db)
	serveJson(ctx, http.StatusOK, &info)
}

func ServiceClear(ctx *fasthttp.RequestCtx) {
	info := models.DatabaseInfo{}
	info.Clean(db)
	serveJson(ctx, http.StatusOK, nil)
}
