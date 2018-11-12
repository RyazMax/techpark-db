package main

// docker build -t a.navrotskiy https://github.com/bozaro/tech-db-forum-server.git
// docker run -p 5000:5000 --name a.navrotskiy -t a.navrotskiy
// ./tech-db-forum func -u http://localhost:5000/api -r report.html

import (
	"log"
	"net/http"
	"techpark-db/database"
	"techpark-db/handlers"

	"github.com/gorilla/mux"
)

func main() {
	var db database.DB
	db.ConectDB()
	db.InitDB("database/init.sql")
	router := mux.NewRouter().PathPrefix("/api").Subrouter()
	router.Handle("/user/{name}/create", &handlers.CreateUserHandler{DB: &db})
	router.Handle("/user/{name}/profile", &handlers.ProfileUserHandler{DB: &db})
	router.Handle("/forum/create", &handlers.CreateForumHandler{DB: &db})
	log.Println("started server")
	http.ListenAndServe(":5000", router)
}
