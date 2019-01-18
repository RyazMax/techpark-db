package models

import (
	"log"
	"techpark-db/database"
)

type DatabaseInfo struct {
	ForumCount  int `json:"forum"`
	PostCount   int `json:"post"`
	ThreadCount int `json:"thread"`
	UserCount   int `json:"user"`
}

func GetInfo(db *database.DB) (d DatabaseInfo) {
	err := db.DataBase.QueryRow(`
		SELECT * FROM
		(( SELECT COUNT(*) FROM forum) AS forums
		CROSS JOIN (SELECT COUNT(*) FROM post) AS posts
		CROSS JOIN (SELECT COUNT(*) FROM thread) AS threads
		CROSS JOIN (SELECT COUNT(*) FROM forum_user) AS users);`).Scan(&d.ForumCount, &d.PostCount, &d.ThreadCount, &d.UserCount)

	if err != nil {
		log.Println(err)
	}
	return
}

func Clean(db *database.DB) {
	_, err := db.DataBase.Exec("TRUNCATE forum_user, forum, thread, post, vote, user_in_forum;")
	if err != nil {
		log.Println(err)
	}
}
