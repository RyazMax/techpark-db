package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

func ThreadVote(ctx *fasthttp.RequestCtx) {
	slugOrID := ctx.UserValue("slug_or_id").(string)
	id, _ := strconv.Atoi(slugOrID)
	body := ctx.PostBody()

	thread := models.Thread{}
	var exist bool
	if id != 0 {
		exist = thread.GetById(id, db)
	} else {
		exist = thread.GetBySlug(slugOrID, db)
	}
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Thread not found"})
		return
	}

	vote := models.Vote{}
	json.Unmarshal(body, &vote)

	author := models.User{}
	exist = author.GetUserByNick(vote.Nickname, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Author not found"})
		return
	}
	vote.Thread = thread.ID
	vote.Nickname = author.Nickname

	vote.Add(db)

	thread.Votes = thread.GetVotesById(thread.ID, db)
	serveJson(ctx, http.StatusOK, &thread)
}
