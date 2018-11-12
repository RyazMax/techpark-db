package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"techpark-db/database"
	"techpark-db/models"
)

type CreateForumHandler struct {
	DB *database.DB
}

func createForum(w http.ResponseWriter, r *http.Request, db *database.DB) {
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Can not read body")
	}

	forum := models.Forum{}
	err = json.Unmarshal(body, &forum)
	if err != nil {
		log.Println(err)
	}

	user := models.User{}
	exist := user.GetUserByNick(forum.User, db)
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		msg := models.Message{Message: "User not found"}
		body, err = json.Marshal(msg)
	} else if forum.GetBySlug(forum.Slug, db) {
		w.WriteHeader(http.StatusConflict)
		msg := models.Message{Message: "Forum already exist"}
		body, err = json.Marshal(msg)
	} else {
		err = forum.Create(db)
		if err != nil {
			log.Println(err)
		}
	}

	err = forum.Create(db)
	w.WriteHeader(http.StatusCreated)
	w.Write(body)
}

func (h *CreateForumHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createForum(w, r, h.DB)
	}
}
