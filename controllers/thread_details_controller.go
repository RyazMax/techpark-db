package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

func ThreadDetailsGet(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.Atoi(slugOrID)

	thread := models.Thread{}
	var exist bool
	if err == nil {
		thread, exist = models.ThreadGetById(id, db)
	} else {
		thread, exist = models.ThreadGetBySlug(slugOrID, db)
	}
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Thread not found"})
		return
	}

	serveJson(ctx, http.StatusOK, &thread)
}

func ThreadDetailsPost(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.Atoi(slugOrID)
	body := ctx.PostBody()

	thread := models.Thread{}
	var exist bool
	if err == nil {
		thread, exist = models.ThreadGetById(id, db)
	} else {
		thread, exist = models.ThreadGetBySlug(slugOrID, db)
	}
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Thread not found"})
		return
	}

	updateThread := models.ThreadUpdate{}
	json.Unmarshal(body, &updateThread)

	if updateThread.Message != "" {
		thread.Message = updateThread.Message
	}
	if updateThread.Title != "" {
		thread.Title = updateThread.Title
	}

	models.ThreadUpd(thread, db)
	serveJson(ctx, http.StatusOK, &thread)
}
