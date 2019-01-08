package controllers

import (
	"net/http"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

type ServiceController struct {
	beego.Controller
	DB *database.DB
}

func (c *ServiceController) Get() {
	c.Ctx.Output.SetStatus(http.StatusOK)
	info := models.DatabaseInfo{}
	info.Get(c.DB)
	c.Data["json"] = &info
	c.ServeJSON()
}

func (c *ServiceController) Post() {
	info := models.DatabaseInfo{}
	info.Clean(c.DB)
	c.Ctx.Output.SetStatus(http.StatusOK)
	c.ServeJSON()
}
