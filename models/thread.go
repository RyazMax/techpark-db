package models

import "techpark-db/database"

type Thread struct {
	Author  string `json:"author"`
	Created string `json:"created"`
	Forum   string `json:"forum"`
	Id      int    `json:"id"`
	Message string `json:"message"`
	Slug    string `json:"slug"`
	Title   string `json:"title"`
	Votes   int    `json:"votes"`
}

func (t *Thread) Add(db *database.DB) {

}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type Threads []Thread
