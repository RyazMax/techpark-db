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
	json.Unmarshal(body, newUser)

	newUser.Nickname = nickname
	sameUsers := newUser.GetLike(c.DB)
	if len(sameUsers) > 0 {
		//beego.Info(sameUsers)
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Data["json"] = &sameUsers
		c.ServeJSON()
		return
	}

	newUser.Add(c.DB)
	c.Ctx.Output.SetStatus(http.StatusCreated)
	c.Data["json"] = &newUser
	c.ServeJSON()
}
