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

	err = json.Unmarshal(body, &posts)
	if err != nil {
		beego.Warn("Can not unmarshal body", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = &models.Message{Message: "Can not unmarshal"}
		c.ServeJSON()
		return
	}

	posts, err = thread.AddPosts(posts, c.DB)
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
	forum := models.Forum{}
	forum.GetBySlug(thread.Forum, c.DB)
	forum.Posts += len(posts)
	forum.Update(c.DB)

	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = &posts
	c.ServeJSON()
}
