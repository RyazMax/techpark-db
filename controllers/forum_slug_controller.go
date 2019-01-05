package controllers

import (
	"net/http"
	"path"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

type ForumSlugController struct {
	beego.Controller
	DB *database.DB
}

func (c *ForumSlugController) Get() {
	switch path.Base(c.Ctx.Input.URL()) {
	case "details":
		c.details()
		return
	case "threads":
		c.threads()
		return
	case "users":
		c.users()
		return
	}
}

func (c *ForumSlugController) details() {
	forum := models.Forum{}
	slug := c.Ctx.Input.Param(":slug")

	exist := forum.GetBySlug(slug, c.DB)
	if exist {
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.Data["json"] = &forum
		c.ServeJSON()
	} else {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Can not find forum"}
		c.ServeJSON()
	}
}

func (c *ForumSlugController) threads() {

}

func (c *ForumSlugController) users() {

}
