package controllers

import (
	"encoding/json"
	"net/http"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

type ForumController struct {
	beego.Controller
	DB *database.DB
}

func (c *ForumController) Post() {
	body := c.Ctx.Input.RequestBody
	newForum := models.Forum{}

	json.Unmarshal(body, &newForum)

	owner := models.User{}
	exist := owner.GetUserByNick(newForum.User, c.DB)
	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Can not find user\n"}
		c.ServeJSON()
		return
	}
	newForum.User = owner.Nickname

	oldForum := models.Forum{}
	exist = oldForum.GetBySlug(newForum.Slug, c.DB)
	if exist {
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Data["json"] = &oldForum
		c.ServeJSON()
		return
	}

	newForum.Create(c.DB)
	newForum.GetBySlug(newForum.Slug, c.DB)
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = &newForum
	c.ServeJSON()
}
