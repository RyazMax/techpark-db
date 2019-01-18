package models

import (
	"techpark-db/database"
)

type Vote struct {
	Nickname string `json:"nickname"`
	Voice    int    `json:"voice"`
	Thread   int    `json:"thread,ommitempty"`
}

func GetVoteByNickAndID(nick string, id int, db *database.DB) (v Vote, exist bool) {
	err := db.DataBase.QueryRow("SELECT * FROM vote WHERE nickname = $1 AND thread=$2;", nick, id).
		Scan(&v.Nickname, &v.Voice, &v.Thread)
	if err != nil {
		return Vote{}, false
	}
	return v, true
}

func VoteAdd(v Vote, db *database.DB) {
	_, err := db.DataBase.Exec("INSERT INTO vote as v (nickname,voice,thread) values($1,$2,$3) ON CONFLICT (nickname, thread) DO UPDATE SET voice=$2;", v.Nickname, v.Voice, v.Thread)
	if err != nil {
		return
	}
}

func VoteUpd(v Vote, db *database.DB) {
	db.DataBase.Exec("UPDATE vote SET voice = $1 WHERE nickname = $2 AND thread=$3;",
		v.Voice, v.Nickname, v.Thread)
}
