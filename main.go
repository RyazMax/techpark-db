package main

// docker build -t a.navrotskiy https://github.com/RyazMax/techpark-db.git
// docker run -p 5000:5000 --name a.navrotskiy -t a.navrotskiy
// ./tech-db-forum func -u http://localhost:5000/api -r report.html

import (
	"log"
	"techpark-db/database"
	"techpark-db/routers"

	"github.com/valyala/fasthttp"
)

func main() {

	db := database.GetDB()
	db.GetPool()
	db.InitDB("database/init.sql")
	db.DataBase.Close()
	db.GetPool()

	//db.InitDB("database/init.sql")
	//defer db.DataBase.Close()
	log.Println("Start server on :5000")
	log.Fatal(fasthttp.ListenAndServe(":5000", routers.Handler))
}
