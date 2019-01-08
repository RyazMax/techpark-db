package controllers

import (
	"encoding/json"
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

func (c *ForumSlugController) Post() {
	slug := c.Ctx.Input.Param(":slug")
	body := c.Ctx.Input.RequestBody

	newThread := models.Thread{}
	err := json.Unmarshal(body, &newThread)
	if err != nil {
		beego.Warn("Can not unmarshal body", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = &models.Message{Message: "Can not unmarshal"}
		c.ServeJSON()
		return
	}

	// Наличие юзера
	owner := models.User{}
	exist := owner.GetUserByNick(newThread.Author, c.DB)
	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Can not find user"}
		c.ServeJSON()
		return
	}
	newThread.Author = owner.Nickname

	// Наличие форума
	forum := models.Forum{}
	exist = forum.GetBySlug(slug, c.DB)
	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Can not find forum"}
		c.ServeJSON()
		return
	}

	// Дубликат по ветке
	newThread.Forum = forum.Slug
	if newThread.Slug != "" {
		oldThread := models.Thread{}
		exist = oldThread.GetBySlug(newThread.Slug, c.DB)
		if exist {
			c.Ctx.Output.SetStatus(http.StatusConflict)
			c.Data["json"] = &oldThread
			c.ServeJSON()
			return
		}
	}

	//forum.Threads++
	//forum.Update(c.DB)
	newThread.Add(c.DB)
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = &newThread
	c.ServeJSON()
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
	var forum models.Forum
	slug := c.Ctx.Input.Param(":slug")
	exist := forum.GetBySlug(slug, c.DB)
	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = models.Message{Message: "Can not find forum"}
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
	threads := models.GetThreadsSorted(slug, limit, since, desc, c.DB)

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = &threads
	c.ServeJSON()
}

func (c *ForumSlugController) users() {
	var forum models.Forum
	slug := c.Ctx.Input.Param(":slug")
	exist := forum.GetBySlug(slug, c.DB)
	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = models.Message{Message: "Can not find forum"}
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
	users := models.GetUsersSorted(slug, limit, since, desc, c.DB)

	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = &users
	c.ServeJSON()
}
