package controllers

import (
	"encoding/json"
	"net/http"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

func ForumCreate(ctx *fasthttp.RequestCtx) {
	body := ctx.PostBody()
	newForum := models.Forum{}

	json.Unmarshal(body, &newForum)

	owner, exist := models.GetUserByNick(newForum.User, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Can not find user"})
		return
	}
	newForum.User = owner.Nickname

	oldForum, exist := models.ForumGetBySlug(newForum.Slug, db)
	if exist {
		serveJson(ctx, http.StatusConflict, &oldForum)
		return
	}

	models.ForumCreate(newForum, db)
	serveJson(ctx, http.StatusCreated, &newForum)

}
