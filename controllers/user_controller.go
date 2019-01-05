package controllers

import (
	"encoding/json"
	"net/http"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

// UserController is main controller for my app
type UserController struct {
	beego.Controller
	DB *database.DB
}

// Post method createas new user
func (c *UserController) Post() {
	nickname := c.Ctx.Input.Param(":nickname")
	body := c.Ctx.Input.RequestBody
	newUser := &models.User{}
	err := json.Unmarshal(body, newUser)
	if err != nil {
		beego.Warn("Can not unmarshal body", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = &models.Message{Message: "Can not unmarshal"}
		c.ServeJSON()
		return
	}

	newUser.Nickname = nickname
	sameUsers := newUser.GetLike(c.DB)
	if len(sameUsers) > 0 {
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Data["json"] = &sameUsers
		c.ServeJSON()
		return
	}

	err = newUser.Add(c.DB)
	if err != nil {
		beego.Error(err)
		return
	}
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = &newUser
	c.ServeJSON()
}
