package controllers

import (
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type ThreadDetailsController struct {
	beego.Controller
	DB *database.DB
}
