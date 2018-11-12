package handlers

import (
	"net/http"
	"techpark-db/database"
)

type PostHanlder struct {
	DB *database.DB
}

func (h *PostHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
