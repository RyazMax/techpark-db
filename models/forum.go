package models

import (
	"log"
	"techpark-db/database"
)

type Forum struct {
	Posts   int    `json:"posts,ommitempty"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads,ommitempty"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

func ForumCreate(f Forum, db *database.DB) error {
	_, err := db.DataBase.Exec("insert into forum(posts,slug,threads,title,forum_user) values ($1,$2,$3,$4,$5);",
		f.Posts, f.Slug, f.Threads, f.Title, f.User)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func ForumUpdate(f Forum, db *database.DB) {
	_, err := db.DataBase.Exec("UPDATE forum SET threads=$1, posts=$2 WHERE slug = $3;", f.Threads, f.Posts, f.Slug)
	if err != nil {
		log.Println(err)
	}
}

func ForumGetBySlug(slug string, db *database.DB) (f Forum, exist bool) {
	err := db.DataBase.QueryRow("select * from forum where slug = $1;", slug).
		Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User)
	if err != nil {
		return Forum{}, false
	}

	return f, true
}
