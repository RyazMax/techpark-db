package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/astaxie/beego"
)

type ThreadVoteController struct {
	beego.Controller
	DB *database.DB
}

func (c *ThreadVoteController) Post() {
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

	vote := models.Vote{}
	err = json.Unmarshal(body, &vote)
	if err != nil {
		beego.Warn("Can not unmarshal body", err)
		c.Ctx.Output.SetStatus(http.StatusBadRequest)
		c.Data["json"] = &models.Message{Message: "Can not unmarshal"}
		c.ServeJSON()
		return
	}

	author := models.User{}
	exist = author.GetUserByNick(vote.Nickname, c.DB)
	if !exist {
		c.Ctx.Output.SetStatus(http.StatusNotFound)
		c.Data["json"] = &models.Message{Message: "Author not found"}
		c.ServeJSON()
		return
	}
	vote.Thread = thread.ID
	vote.Nickname = author.Nickname

	//oldVote := models.Vote{}
	//exist = oldVote.GetByNickAndID(vote.Nickname, thread.ID, c.DB)
	//delta := 0
	/*if exist {
		delta = vote.Voice - oldVote.Voice
	} else {
		delta = vote.Voice
	}
	thread.Votes += delta*/

	/*if exist {
		vote.Update(c.DB)
	} else {
		vote.Add(c.DB)
	}*/
	err = vote.Add(c.DB)
	//pqErr := err.(*pq.Error)
	//beego.Info(pqErr.Code)
	thread.GetById(thread.ID, c.DB)
	c.Data["json"] = &thread
	c.ServeJSON()
}
