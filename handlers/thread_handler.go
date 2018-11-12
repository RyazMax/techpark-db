package handlers

import (
	"net/http"
	"techpark-db/database"
)

type ThreadHanlder struct {
	DB *database.DB
}

func (h *ThreadHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
