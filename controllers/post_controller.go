package controllers

import (
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type PostController struct {
	beego.Controller
	DB *database.DB
}
