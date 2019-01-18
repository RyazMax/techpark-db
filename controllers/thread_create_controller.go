package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

func ThreadCreatePosts(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.Atoi(slugOrID)
	body := ctx.PostBody()

	thread := models.Thread{}
	var exist bool
	if id != 0 {
		exist = thread.GetById(id, db)
	} else {
		exist = thread.GetBySlug(slugOrID, db)
	}
	if !exist {
		serveJson(ctx, http.StatusNotFound, models.Message{Message: "Thread not found"})
		return
	}

	posts := models.Posts{}

	json.Unmarshal(body, &posts)

	ids, curTime, err := thread.AddPosts(posts, db)
	if err != nil && err.Error() == "No author" {
		serveJson(ctx, http.StatusNotFound, models.Message{Message: "Author not found"})
		return
	} else if err != nil {
		serveJson(ctx, http.StatusConflict, &models.Message{Message: err.Error()})
		return
	}
	for i := range posts {
		posts[i].Id = ids[i]
		posts[i].Created = curTime
		posts[i].Thread = thread.ID
		posts[i].Forum = thread.Forum
	}

	serveJson(ctx, http.StatusCreated, &posts)
}
