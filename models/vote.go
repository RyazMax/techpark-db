package models

import (
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   int    `json:"thread,ommitempty"`
}

func (v *Vote) GetByNickAndID(nick string, id int, db *database.DB) bool {
	err := db.DataBase.QueryRow("SELECT * FROM vote WHERE nickname = $1 AND thread=$2;", nick, id).
		Scan(&v.Nickname, &v.Voice, &v.Thread)
	if err != nil {
		return false
	}
	return true
}

func (v *Vote) Add(db *database.DB) {
	_, err := db.DataBase.Exec("INSERT INTO vote(nickname,voice,thread) values($1,$2,$3);", v.Nickname, v.Voice, v.Thread)
	if err != nil {
		beego.Warn(err)
	}
}

func (v *Vote) Update(db *database.DB) {
	_, err := db.DataBase.Exec("UPDATE vote SET voice = $1 WHERE nickname = $2 AND thread=$3",
		v.Voice, v.Nickname, v.Thread)
	if err != nil {
		beego.Warn(err)
	}
}
