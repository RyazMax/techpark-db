package models

import (
	"log"
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type Forum struct {
	Posts   int    `json:"posts,ommitempty"`
	Slug    string `json:"slug"`
	Threads int    `json:"threads,ommitempty"`
	Title   string `json:"title"`
	User    string `json:"user"`
}

func (f *Forum) Create(db *database.DB) error {
	_, err := db.DataBase.Exec("insert into forum(posts,slug,threads,title,forum_user) values ($1,$2,$3,$4,$5);",
		f.Posts, f.Slug, f.Threads, f.Title, f.User)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (f *Forum) Update(db *database.DB) {
	_, err := db.DataBase.Exec("UPDATE forum SET threads=$1, posts=$2 WHERE LOWER(slug)=LOWER($3);", f.Threads, f.Posts, f.Slug)
	if err != nil {
		beego.Warn(err)
	}
}

func (f *Forum) GetBySlug(slug string, db *database.DB) (exist bool) {
	rows, err := db.DataBase.Query("select * from forum where LOWER(slug)=LOWER($1);", slug)
	defer rows.Close()
	if err != nil {
		log.Println(err)
		return false
	}

	for rows.Next() {
		exist = true
		err = rows.Scan(&f.Posts, &f.Slug, &f.Threads, &f.Title, &f.User)
		if err != nil {
			log.Println(err)
			return false
		}
	}
	return
}
