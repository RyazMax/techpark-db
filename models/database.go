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
	rows, err := db.DataBase.Query("SELECT COUNT(*) FROM forum; SELECT COUNT(*) FROM post; SELECT COUNT(*) FROM thread; SELECT COUNT(*) FROM forum_user;")
	defer rows.Close()
	if err != nil {
		beego.Warn(err)
	}

	rows.Next()
	err = rows.Scan(&d.ForumCount)
	rows.Next()
	rows.Next()
	err = rows.Scan(&d.PostCount)
	rows.Next()
	rows.Next()
	err = rows.Scan(&d.ThreadCount)
	rows.Next()
	rows.Next()
	err = rows.Scan(&d.UserCount)

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
