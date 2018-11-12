package handlers

import (
	"net/http"
	"techpark-db/database"
)

type VoteHanlder struct {
	DB *database.DB
}

func (h *VoteHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
