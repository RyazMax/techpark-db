package controllers

import (
	"encoding/json"
	"net/http"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

var db *database.DB

func init() {
	db = database.GetDB()
}

func ForumSlugCreate(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	body := ctx.PostBody()

	newThread := models.Thread{}
	json.Unmarshal(body, &newThread)

	// Наличие юзера
	owner, exist := models.GetUserByNick(newThread.Author, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Can not find user"})
		return
	}
	newThread.Author = owner.Nickname

	// Наличие форума
	forum, exist := models.ForumGetBySlug(slug, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Can not find forum"})
		return
	}

	// Дубликат по ветке
	newThread.Forum = forum.Slug
	if newThread.Slug != "" {
		oldThread, exist := models.ThreadGetBySlug(newThread.Slug, db)
		if exist {
			serveJson(ctx, http.StatusConflict, &oldThread)
			return
		}
	}

	models.ThreadAdd(&newThread, db)
	serveJson(ctx, http.StatusCreated, &newThread)
}

func ForumSlugDetails(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	forum, exist := models.ForumGetBySlug(slug, db)
	if exist {
		serveJson(ctx, http.StatusOK, &forum)
	} else {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Can not find forum"})
	}
}

func ForumSlugThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	_, exist := models.ForumGetBySlug(slug, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Can not find forum"})
		return
	}

	var desc bool
	descString := string(ctx.QueryArgs().Peek("desc"))
	if descString == "true" {
		desc = true
	}
	limit := ctx.QueryArgs().GetUintOrZero("limit")
	since := string(ctx.QueryArgs().Peek("since"))
	threads := models.GetThreadsSorted(slug, limit, since, desc, db)

	serveJson(ctx, http.StatusOK, &threads)
}

func ForumSlugUsers(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)
	_, exist := models.ForumGetBySlug(slug, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Can not find forum"})
		return
	}

	limit := ctx.QueryArgs().GetUintOrZero("limit")

	var desc bool
	descString := string(ctx.QueryArgs().Peek("desc"))
	if descString == "true" {
		desc = true
	}

	since := string(ctx.QueryArgs().Peek("since"))
	users := models.GetUsersSorted(slug, limit, since, desc, db)

	serveJson(ctx, http.StatusOK, &users)
}
