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
		thread, exist = models.ThreadGetById(id, db)
	} else {
		thread, exist = models.ThreadGetBySlug(slugOrID, db)
	}
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Thread not found"})
		return
	}

	vote := models.Vote{}
	json.Unmarshal(body, &vote)

	author, exist := models.GetUserByNick(vote.Nickname, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Author not found"})
		return
	}
	vote.Thread = thread.ID
	vote.Nickname = author.Nickname

	models.VoteAdd(vote, db)

	thread.Votes = models.GetVotesById(thread.ID, db)
	serveJson(ctx, http.StatusOK, &thread)
}
