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
	post, exist := models.PostGetByID(id, db)
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
		models.PostUpd(post, db)
	}

	serveJson(ctx, http.StatusOK, &post)
}

func PostDetailsGet(ctx *fasthttp.RequestCtx) {
	param := ctx.UserValue("id").(string)
	id, _ := strconv.Atoi(param)
	post, exist := models.PostGetByID(id, db)
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
		author, _ := models.GetUserByNick(post.Author, db)
		result.Author = &author
	}

	if withThread {
		thread, _ := models.ThreadGetById(post.Thread, db)
		result.Thread = &thread
	}

	if withForum {
		forum, _ := models.ForumGetBySlug(post.Forum, db)
		result.Forum = &forum
	}

	serveJson(ctx, http.StatusOK, &result)
}
