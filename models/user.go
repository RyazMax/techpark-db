package models

import (
	"fmt"
	"log"
	"strings"
	"techpark-db/database"

	"github.com/astaxie/beego"
	"github.com/jackc/pgx"
)

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
}

type UserUpdate struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

//easyjson:json
type Users []User

func UserAdd(newUser User, db *database.DB) error {
	_, err := db.DataBase.Exec("insert into forum_user(email,about,fullname,nickname) values ($1,$2,$3,$4);", newUser.Email, newUser.About, newUser.Fullname, newUser.Nickname)
	if err != nil {
		return err
	}
	return nil
}

func AddUsersToForum(forum string, users *map[string]bool, db *database.DB) {
	var query strings.Builder
	query.Grow(80 + 11*len(*users))
	query.WriteString("INSERT into user_in_forum(nickname, forum) VALUES ")
	var cnt int
	args := make([]interface{}, 1, len(*users))
	args[0] = forum
	for nick := range *users {
		if cnt > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("($%d, $%d)", cnt+2, 1))
		args = append(args, nick)
		cnt++
	}
	query.WriteString("ON CONFLICT DO NOTHING;")
	db.DataBase.Exec(query.String(), args...)
}

func GetUserByNickOrEmail(nick string, email string, db *database.DB) Users {
	rows, err := db.DataBase.Query("select * from forum_user where nickname = $2 or email = $1;",
		email, nick)
	defer rows.Close()
	users := make(Users, 0, 2)
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Email, &user.About, &user.Fullname, &user.Nickname)
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, user)
	}

	return users
}

func GetUserByNick(nickname string, db *database.DB) (u User, exist bool) {
	err := db.DataBase.QueryRow("select * from forum_user where nickname = $1;", nickname).
		Scan(&u.Email, &u.About, &u.Fullname, &u.Nickname)
	if err != nil {
		return User{}, false
	}
	return u, true
}

func UserUpd(u User, db *database.DB) error {

	var err error
	if u.About != "" {
		_, err = db.DataBase.Exec("UPDATE forum_user SET about = $1 WHERE nickname = $2", u.About, u.Nickname)
	}
	if u.Fullname != "" {
		_, err = db.DataBase.Exec("UPDATE forum_user SET fullname = $1 WHERE nickname =  $2", u.Fullname, u.Nickname)
	}
	if u.Email != "" {
		_, err = db.DataBase.Exec("UPDATE forum_user SET email = $1 WHERE nickname = $2", u.Email, u.Nickname)
	}
	//_, err := db.DataBase.Exec("UPDATE forum_user SET fullname = $1, about = $2, email = $3 WHERE LOWER(nickname)=LOWER($4);", u.Fullname, u.About, u.Email, u.Nickname)
	if err != nil {
		log.Println("Update", err)
	}
	return err
}

func GetUsersSorted(slug string, limit int, since string, desc bool, db *database.DB) Users {
	var (
		rows *pgx.Rows
		err  error
	)
	var query strings.Builder
	query.WriteString(`select u.* from forum_user u JOIN user_in_forum uf ON u.nickname=uf.nickname where forum=$1 `)
	if since != "" {
		query.WriteString("AND u.nickname ")
		if desc {
			query.WriteString("< $2 ")
		} else {
			query.WriteString("> $2 ")
		}
		query.WriteString("GROUP BY u.nickname ORDER BY u.nickname ")
		if desc {
			query.WriteString("DESC ")
		}
		if limit != 0 {
			query.WriteString("LIMIT $3;")
			rows, err = db.DataBase.Query(query.String(), slug, since, limit)
		} else {
			query.WriteString(";")
			rows, err = db.DataBase.Query(query.String(), slug, since)
		}
	} else {
		query.WriteString("ORDER BY u.nickname ")
		if desc {
			query.WriteString("DESC ")
		}
		if limit != 0 {
			query.WriteString("LIMIT $2;")
			rows, err = db.DataBase.Query(query.String(), slug, limit)
		} else {
			query.WriteString(";")
			rows, err = db.DataBase.Query(query.String(), slug)
		}
	}
	defer rows.Close()

	users := make(Users, 0)
	for rows.Next() {
		var u User
		err = rows.Scan(&u.Email, &u.About, &u.Fullname, &u.Nickname)
		if err != nil {
			return nil
		}
		users = append(users, u)
	}
	return users
}

func GetUsersByNicks(nicks *map[string]bool, db *database.DB) (res int) {
	if len(*nicks) == 0 {
		return 0
	}
	var query strings.Builder
	args := make([]interface{}, 0, len(*nicks))
	query.Grow(100 + len(*nicks)*2)
	query.WriteString("SELECT nickname FROM forum_user WHERE nickname in (")
	var cnt int
	for name := range *nicks {
		if cnt > 0 {
			query.WriteString(", ")
		}
		cnt++
		query.WriteString(fmt.Sprintf("$%d ", cnt))
		args = append(args, name)
	}
	query.WriteString(");")

	rows, err := db.DataBase.Query(query.String(), args...)
	defer rows.Close()
	if err != nil {
		beego.Warn(err)
		return 0
	}
	var tmp string
	for rows.Next() {
		rows.Scan(&tmp)
		res++
	}
	return res
}
