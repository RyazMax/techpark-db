package controllers

import (
	"net/http"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

func ServiceStatus(ctx *fasthttp.RequestCtx) {
	info := models.GetInfo(db)
	serveJson(ctx, http.StatusOK, &info)
}

func ServiceClear(ctx *fasthttp.RequestCtx) {
	models.Clean(db)
	serveJson(ctx, http.StatusOK, nil)
}
