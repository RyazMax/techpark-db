package models

import "techpark-db/database"

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
}

func (v *Vote) Add(db *database.DB) {

}
