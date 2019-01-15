package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

type ThreadCreateController struct {
	beego.Controller
	DB *database.DB
}

func (c *ThreadCreateController) Post() {
	slugOrID := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(slugOrID)
	body := c.Ctx.Input.RequestBody

	thread := models.Thread{}
	var exist bool
	if id != 0 {
		exist = thread.GetById(id, c.DB)
	} else {
		exist = thread.GetBySlug(slugOrID, c.DB)
	}
	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Thread not found"}
		c.ServeJSON()
		return
	}

	posts := models.Posts{}

	json.Unmarshal(body, &posts)

	ids, curTime, err := thread.AddPosts(posts, c.DB)
	//beego.Info("IN CONTROLL posts_len ", len(posts))
	if err != nil && err.Error() == "No author" {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Author not found"}
		c.ServeJSON()
		return
	} else if err != nil {
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Data["json"] = &models.Message{Message: err.Error()}
		c.ServeJSON()
		return
	}
	for i, _ := range posts {
		posts[i].Id = ids[i]
		posts[i].Created = curTime
		posts[i].Thread = thread.ID
		posts[i].Forum = thread.Forum
	}

	c.Ctx.Output.SetStatus(http.StatusCreated)
	//beego.Info("OK")
	c.Data["json"] = &posts
	c.ServeJSON()
}
