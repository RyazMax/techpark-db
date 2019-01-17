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

func (newUser *User) Add(db *database.DB) error {
	_, err := db.DataBase.Exec("insert into forum_user(email,about,fullname,nickname) values ($1,$2,$3,$4);", newUser.Email, newUser.About, newUser.Fullname, newUser.Nickname)
	if err != nil {
		beego.Warn("IN ADD", err)
		return err
	}
	return nil
}

func AddUsersToForum(forum string, users *map[string]bool, db *database.DB) {
	var query strings.Builder
	query.WriteString("INSERT into user_in_forum(nickname, forum) VALUES ")
	var cnt int
	for nick := range *users {
		if cnt > 0 {
			query.WriteString(", ")
		}
		query.WriteString(fmt.Sprintf("('%s', '%s')", nick, forum))
	}
	query.WriteString("ON CONFLICT DO NOTHING;")
	db.DataBase.Exec(query.String())
}

func (u *User) GetLike(db *database.DB) Users {
	rows, err := db.DataBase.Query("select * from forum_user where nickname = $2 or email = $1;",
		u.Email, u.Nickname)
	if err != nil {
		beego.Warn("IN get like ", err.Error())
		//beego.Warn(err.(*pgx.PgError).Detail)
		//beego.Warn(err.(*pgx.PgError).Message)
	}
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

func (u *User) GetUserByNick(nickname string, db *database.DB) (exist bool) {
	err := db.DataBase.QueryRow("select * from forum_user where nickname = $1;", nickname).
		Scan(&u.Email, &u.About, &u.Fullname, &u.Nickname)
	if err != nil {
		return false
	}
	return true
}

func (u *User) Update(db *database.DB) error {

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
	// Нет вложенного
	/*subQuery :=
	`select u.* from forum_user u
	JOIN
	((select distinct author from thread WHERE forum = $1)
	UNION
	(select distinct author from post WHERE forum = $1)) as p ON nickname=p.author `*/
	subQuery := `select u.* from forum_user u JOIN user_in_forum uf ON u.nickname=uf.nickname where forum=$1 `
	if since != "" {
		subQuery += "AND u.nickname "
		if desc {
			subQuery += "< $2 "
		} else {
			subQuery += "> $2 "
		}
		subQuery += "GROUP BY u.nickname ORDER BY u.nickname "
		if desc {
			subQuery += "DESC "
		}
		if limit != 0 {
			rows, err = db.DataBase.Query(subQuery+"LIMIT $3;", slug, since, limit)
		} else {
			rows, err = db.DataBase.Query(subQuery+";", slug, since)
		}
	} else {
		subQuery += "ORDER BY u.nickname "
		if desc {
			subQuery += "DESC "
		}
		if limit != 0 {
			rows, err = db.DataBase.Query(subQuery+"LIMIT $2;", slug, limit)
		} else {
			rows, err = db.DataBase.Query(subQuery, slug)
		}
	}
	defer rows.Close()
	if err != nil {
		beego.Warn(err)
		return nil
	}
	users := make(Users, 0)
	for rows.Next() {
		var u User
		err = rows.Scan(&u.Email, &u.About, &u.Fullname, &u.Nickname)
		if err != nil {
			beego.Warn(err)
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
	query.WriteString("SELECT COUNT (*) FROM forum_user WHERE nickname in (")
	var cnt int
	for name := range *nicks {
		if cnt > 0 {
			query.WriteString(", ")
		}
		cnt++
		query.WriteString(fmt.Sprintf("'%s' ", name))
	}
	query.WriteString(");")

	err := db.DataBase.QueryRow(query.String()).Scan(&res)
	if err != nil {
		beego.Warn(err)
		return 0
	}
	return res
}
