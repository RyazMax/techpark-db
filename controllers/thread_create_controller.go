package controllers

import (
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type ThreadCreateController struct {
	beego.Controller
	DB *database.DB
}
