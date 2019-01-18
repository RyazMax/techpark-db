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

	owner := models.User{}
	exist := owner.GetUserByNick(newForum.User, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Can not find user"})
		return
	}
	newForum.User = owner.Nickname

	oldForum := models.Forum{}
	exist = oldForum.GetBySlug(newForum.Slug, db)
	if exist {
		serveJson(ctx, http.StatusConflict, &oldForum)
		return
	}

	newForum.Create(db)
	newForum.GetBySlug(newForum.Slug, db)
	serveJson(ctx, http.StatusCreated, &newForum)

}
