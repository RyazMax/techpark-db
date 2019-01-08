package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

type PostController struct {
	beego.Controller
	DB *database.DB
}

func (c *PostController) Post() {
	param := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(param)
	post := models.Post{}
	exist := post.GetByID(id, c.DB)
	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Post not found"}
		c.ServeJSON()
		return
	}

	body := c.Ctx.Input.RequestBody
	postUpdate := models.PostUpdate{}
	err := json.Unmarshal(body, &postUpdate)
	if err != nil {
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = &models.Message{Message: "Can not unmarshal"}
		c.ServeJSON()
		return
	}

	if postUpdate.Message != post.Message && postUpdate.Message != "" {
		post.Message = postUpdate.Message
		post.IsEdited = true
		post.Update(c.DB)
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = &post
	c.ServeJSON()
}

func (c *PostController) Get() {
	param := c.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(param)
	post := models.Post{}
	exist := post.GetByID(id, c.DB)
	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Post not found"}
		c.ServeJSON()
		return
	}

	related := c.GetString("related")
	details := strings.Split(related, ",")
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
		author.GetUserByNick(post.Author, c.DB)
		result.Author = &author
	}

	if withThread {
		thread := models.Thread{}
		thread.GetById(post.Thread, c.DB)
		result.Thread = &thread
	}

	if withForum {
		forum := models.Forum{}
		forum.GetBySlug(post.Forum, c.DB)
		result.Forum = &forum
	}

	c.Data["json"] = &result
	c.ServeJSON()
}
