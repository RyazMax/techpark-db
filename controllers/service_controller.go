package controllers

import (
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type ServiceController struct {
	beego.Controller
	DB *database.DB
}
