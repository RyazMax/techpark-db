package controllers

import (
	"net/http"
	"strconv"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

func ThreadGetPosts(ctx *fasthttp.RequestCtx) {
	var thread models.Thread
	slugOrID := ctx.UserValue("slug_or_id").(string)
	id, err := strconv.Atoi(slugOrID)
	var exist bool
	if err != nil {
		thread, exist = models.ThreadGetBySlug(slugOrID, db)
	} else {
		thread, exist = models.ThreadGetById(id, db)
	}

	if !exist {
		serveJson(ctx, http.StatusNotFound, models.Message{Message: "Can not find thread"})
		return
	}

	limit := ctx.QueryArgs().GetUintOrZero("limit")
	var desc bool
	descString := string(ctx.QueryArgs().Peek("desc"))
	if descString == "true" {
		desc = true
	}
	since := string(ctx.QueryArgs().Peek("since"))
	sortType := string(ctx.QueryArgs().Peek("sort"))
	if sortType == "" {
		sortType = "flat"
	}

	posts := models.Posts{}
	switch sortType {
	case "flat":
		posts = models.GetPostsSortedFlat(thread.ID, limit, since, desc, db)
	case "tree":
		posts = models.GetPostsSortedTree(thread.ID, limit, since, desc, db)
	case "parent_tree":
		posts = models.GetPostsSortedParentTree(thread.ID, limit, since, desc, db)
	}

	serveJson(ctx, http.StatusOK, &posts)
}
