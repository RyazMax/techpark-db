package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

type ThreadDetailsController struct {
	beego.Controller
	DB *database.DB
}

func (c *ThreadDetailsController) Get() {
	slugOrID := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(slugOrID)

	thread := models.Thread{}
	var exist bool
	if err == nil {
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
	c.Data["json"] = &thread
	c.ServeJSON()
}

func (c *ThreadDetailsController) Post() {
	slugOrID := c.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(slugOrID)
	body := c.Ctx.Input.RequestBody

	thread := models.Thread{}
	var exist bool
	if err == nil {
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

	updateThread := models.ThreadUpdate{}
	json.Unmarshal(body, &updateThread)

	if updateThread.Message != "" {
		thread.Message = updateThread.Message
	}
	if updateThread.Title != "" {
		thread.Title = updateThread.Title
	}

	thread.Update(c.DB)
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.Data["json"] = &thread
	c.ServeJSON()
}
