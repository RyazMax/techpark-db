package main

// docker build -t a.navrotskiy https://github.com/bozaro/tech-db-forum-server.git
// docker run -p 5000:5000 --name a.navrotskiy -t a.navrotskiy
// ./tech-db-forum func -u http://localhost:5000/api -r report.html

import (
	"techpark-db/database"
	"techpark-db/routers"

	"github.com/astaxie/beego"
	_ "github.com/gorilla/mux"
)

func main() {

	var db database.DB
	db.GetPool()
	db.InitDB("database/init.sql")
	defer db.DataBase.Close()

	routers.InitRouter(&db)
	beego.Run(":5000")
}
