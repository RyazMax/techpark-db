package controllers

import (
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type ThreadPostsController struct {
	beego.Controller
	DB *database.DB
}
