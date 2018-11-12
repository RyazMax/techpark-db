package models

import (
	"log"
	"techpark-db/database"
)

type Post struct {
	Author   string `json:"author"`
	Created  string `json:"created"`
	Forum    string `json:"forum"`
	Id       int    `json:"id"`
	IsEdited bool   `json:"isEdited"`
	Message  string `json:"message"`
	Parent   int    `json:"parent"`
	Thread   int    `json:"thread"`
}

func (p *Post) Add(db *database.DB) error {
	_, err := db.DataBase.Exec("insert into post(author,created,forum,id,isEdited,message,parent,thread)"+
		"values ($1,$2,$3,$4,$5,$6,$7,$8)",
		p.Author, p.Created, p.Forum, p.Id, p.IsEdited, p.Message, p.Parent, p.Thread)
	if err != nil {
		log.Println(err)
	}
	return err
}

type PostFull struct {
	Author User
	Forum  Forum
	Post   Post
	Thread Thread
}

type PostUpdate struct {
	Message string `json:"message"`
}

type Posts []Post
