package controllers

import (
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type ThreadVoteController struct {
	beego.Controller
	DB *database.DB
}
