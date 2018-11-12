package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"techpark-db/database"
	"techpark-db/models"

	"github.com/gorilla/mux"
)

type CreateUserHandler struct {
	DB *database.DB
}

func createUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
	vars := mux.Vars(r)
	name := vars["name"]
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Can not read body", err)
		return
	}
	newUser := &models.User{}
	err = json.Unmarshal(body, newUser)
	if err != nil {
		log.Println("Can not unmarshal json", err)
		return
	}

	newUser.Nickname = name
	users := newUser.GetLike(db)
	var data []byte
	if len(users) > 0 {
		w.WriteHeader(http.StatusConflict)
		users := newUser.GetLike(db)
		log.Println(users)
		data, err = json.Marshal(users)
		if err != nil {
			log.Println(err)
		}
	} else {
		err = newUser.Add(db)
		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusCreated)
		data, err = json.Marshal(newUser)
		if err != nil {
			log.Println(err)
		}
	}

	log.Println(w.Header())
	w.Write(data)
}

func (h *CreateUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r, h.DB)
	}
}

type ProfileUserHandler struct {
	DB *database.DB
}

func updateUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
	vars := mux.Vars(r)
	name := vars["name"]
	w.Header().Set("Content-Type", "application/json")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Can not read body", err)
		return
	}
	newUser := &models.User{}
	err = json.Unmarshal(body, newUser)
	if err != nil {
		log.Println("Can not unmarshal json", err)
		return
	}

	newUser.Nickname = name
	users := newUser.GetLike(db)
	var data []byte

	if len(users) == 0 || (len(users) == 1 && users[0].Nickname != newUser.Nickname) {
		w.WriteHeader(http.StatusNotFound)
		msg := models.Message{Message: "Error not found"}
		data, err = json.Marshal(msg)
		if err != nil {
			log.Println(err)
		}
	} else if len(users) == 1 {
		w.WriteHeader(http.StatusOK)
		newUser.Update(db)
		data, err = json.Marshal(newUser)
		if err != nil {
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusConflict)
		msg := models.Message{Message: "Error not found"}
		data, err = json.Marshal(msg)
		if err != nil {
			log.Println(err)
		}
	}

	w.Write(data)
}

func getUser(w http.ResponseWriter, r *http.Request, db *database.DB) {
	vars := mux.Vars(r)
	nickname := vars["name"]
	w.Header().Set("Content-Type", "application/json")

	var user models.User

	exist := user.GetUserByNick(nickname, db)

	if exist {
		data, err := json.Marshal(user)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	} else {
		message := models.Message{Message: "Can't find user with id #42\n"}
		data, err := json.Marshal(message)
		if err != nil {
			log.Fatal(err)
		}
		w.WriteHeader(http.StatusNotFound)
		w.Write(data)
	}
}

func (h *ProfileUserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		updateUser(w, r, h.DB)
	case http.MethodGet:
		getUser(w, r, h.DB)
	}
}
