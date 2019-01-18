package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"techpark-db/models"

	"github.com/valyala/fasthttp"
)

func PostDetailsPost(ctx *fasthttp.RequestCtx) {
	param := ctx.UserValue("id").(string)
	id, _ := strconv.Atoi(param)
	post := models.Post{}
	exist := post.GetByID(id, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Post not found"})
		return
	}

	body := ctx.PostBody()
	postUpdate := models.PostUpdate{}
	json.Unmarshal(body, &postUpdate)

	if postUpdate.Message != post.Message && postUpdate.Message != "" {
		post.Message = postUpdate.Message
		post.IsEdited = true
		post.Update(db)
	}

	serveJson(ctx, http.StatusOK, &post)
}

func PostDetailsGet(ctx *fasthttp.RequestCtx) {
	param := ctx.UserValue("id").(string)
	id, _ := strconv.Atoi(param)
	post := models.Post{}
	exist := post.GetByID(id, db)
	if !exist {
		serveJson(ctx, http.StatusNotFound, &models.Message{Message: "Post not found"})
		return
	}

	details := strings.Split(string(ctx.QueryArgs().Peek("related")), ",")
	var (
		withAuthor bool
		withThread bool
		withForum  bool
	)
	if len(details) > 0 {
		for _, one := range details {
			if one == "user" {
				withAuthor = true
			}
			if one == "thread" {
				withThread = true
			}
			if one == "forum" {
				withForum = true
			}
		}
	}

	result := models.PostFull{Post: &post}
	if withAuthor {
		author := models.User{}
		author.GetUserByNick(post.Author, db)
		result.Author = &author
	}

	if withThread {
		thread := models.Thread{}
		thread.GetById(post.Thread, db)
		result.Thread = &thread
	}

	if withForum {
		forum := models.Forum{}
		forum.GetBySlug(post.Forum, db)
		result.Forum = &forum
	}

	serveJson(ctx, http.StatusOK, &result)
}
