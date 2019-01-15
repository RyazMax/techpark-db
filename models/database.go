package models

import (
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type DatabaseInfo struct {
	ForumCount  int `json:"forum"`
	PostCount   int `json:"post"`
	ThreadCount int `json:"thread"`
	UserCount   int `json:"user"`
}

func (d *DatabaseInfo) Get(db *database.DB) {
	err := db.DataBase.QueryRow(`
		SELECT * FROM
		(( SELECT COUNT(*) FROM forum) AS forums
		CROSS JOIN (SELECT COUNT(*) FROM post) AS posts
		CROSS JOIN (SELECT COUNT(*) FROM thread) AS threads
		CROSS JOIN (SELECT COUNT(*) FROM forum_user) AS users);`).Scan(&d.ForumCount, &d.PostCount, &d.ThreadCount, &d.UserCount)

	if err != nil {
		beego.Warn(err)
	}
}

func (d *DatabaseInfo) Clean(db *database.DB) {
	_, err := db.DataBase.Exec("TRUNCATE forum_user, forum, thread, post, vote, user_in_forum;")
	if err != nil {
		beego.Warn(err)
	}
}
