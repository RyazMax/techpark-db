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
		exist = thread.GetBySlug(slugOrID, db)
	} else {
		exist = thread.GetById(id, db)
	}

	if !exist {
		serveJson(ctx, http.StatusNotFound, models.Message{Message: "Can not find thread"})
		return
	}

	limit := ctx.QueryArgs().GetUintOrZero("limit")
	desc := ctx.QueryArgs().GetBool("desc")
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
