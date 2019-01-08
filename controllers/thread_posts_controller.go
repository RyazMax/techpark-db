package controllers

import (
	"net/http"
	"strconv"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

type ThreadPostsController struct {
	beego.Controller
	DB *database.DB
}

func (c *ThreadPostsController) Get() {
	var thread models.Thread
	slugOrID := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(slugOrID)
	var exist bool
	if err != nil {
		exist = thread.GetBySlug(slugOrID, c.DB)
	} else {
		exist = thread.GetById(id, c.DB)
	}

	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = models.Message{Message: "Can not find thread"}
		c.ServeJSON()
		return
	}

	limit, err := c.GetInt("limit", 0)
	if err != nil {
		beego.Warn(err)
	}
	desc, err := c.GetBool("desc", false)
	if err != nil {
		beego.Warn(err)
	}
	since := c.GetString("since", "")
	sortType := c.GetString("sort", "flat")

	posts := models.Posts{}
	switch sortType {
	case "flat":
		posts = models.GetPostsSortedFlat(thread.ID, limit, since, desc, c.DB)
	case "tree":
		posts = models.GetPostsSortedTree(thread.ID, limit, since, desc, c.DB)
	case "parent_tree":
		posts = models.GetPostsSortedParentTree(thread.ID, limit, since, desc, c.DB)
	}

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = &posts
	c.ServeJSON()
}
