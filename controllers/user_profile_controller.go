package controllers

import (
	"encoding/json"
	"net/http"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

// UserProfileController for requests for profile
type UserProfileController struct {
	beego.Controller
	DB *database.DB
}

// Get method returns information about user
func (c *UserProfileController) Get() {
	nickname := c.Ctx.Input.Param(":nickname")

	user := models.User{}
	exist := user.GetUserByNick(nickname, c.DB)

	if exist {
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.Data["json"] = &user
		c.ServeJSON()
	} else {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Can not found user"}
		c.ServeJSON()
	}
}

// Post method returns updates information about user
func (c *UserProfileController) Post() {
	nickname := c.Ctx.Input.Param(":nickname")
	body := c.Ctx.Input.RequestBody

	newUser := models.User{}
	err := json.Unmarshal(body, &newUser)
	if err != nil {
		beego.Warn("Can not unmarshal body", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = &models.Message{Message: "Can not unmarshal"}
		c.ServeJSON()
		return
	}

	newUser.Nickname = nickname
	sameUsers := newUser.GetLike(c.DB)

	if len(sameUsers) == 0 || (len(sameUsers) == 1 && sameUsers[0].Nickname != newUser.Nickname) {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "User not found"}
		c.ServeJSON()
	} else if len(sameUsers) == 1 {
		newUser.Update(c.DB)
		newUser.GetUserByNick(nickname, c.DB)
		c.Ctx.Output.SetStatus(http.StatusOK)
		c.Data["json"] = &newUser
		c.ServeJSON()
	} else {
		c.Ctx.Output.SetStatus(http.StatusConflict)
		c.Data["json"] = &models.Message{Message: "Can not update user"}
		c.ServeJSON()
	}

}
