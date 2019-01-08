package models

import (
	"database/sql"
	"log"
	"techpark-db/database"

	"github.com/astaxie/beego"
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

type Users []User

func (newUser *User) Add(db *database.DB) error {
	_, err := db.DataBase.Exec("insert into forum_user(email,about,fullname,nickname) values ($1,$2,$3,$4);", newUser.Email, newUser.About, newUser.Fullname, newUser.Nickname)
	if err != nil {

		return err
	}
	return nil
}

func (u *User) GetLike(db *database.DB) Users {
	rows, err := db.DataBase.Query("select * from forum_user where LOWER(email)=LOWER($1) or LOWER(nickname)=LOWER($2);",
		u.Email, u.Nickname)
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
	rows, err := db.DataBase.Query("select * from forum_user where LOWER(nickname)=LOWER($1);", nickname)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		exist = true
		err = rows.Scan(&u.Email, &u.About, &u.Fullname, &u.Nickname)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

func (u *User) Update(db *database.DB) error {

	var err error
	if u.About != "" {
		_, err = db.DataBase.Exec("UPDATE forum_user SET about = $1 WHERE LOWER(nickname)=LOWER($2)", u.About, u.Nickname)
	}
	if u.Fullname != "" {
		_, err = db.DataBase.Exec("UPDATE forum_user SET fullname = $1 WHERE LOWER(nickname)=LOWER($2)", u.Fullname, u.Nickname)
	}
	if u.Email != "" {
		_, err = db.DataBase.Exec("UPDATE forum_user SET email = $1 WHERE LOWER(nickname)=LOWER($2)", u.Email, u.Nickname)
	}
	//_, err := db.DataBase.Exec("UPDATE forum_user SET fullname = $1, about = $2, email = $3 WHERE LOWER(nickname)=LOWER($4);", u.Fullname, u.About, u.Email, u.Nickname)
	if err != nil {
		log.Println("Update", err)
	}
	return err
}

func GetUsersSorted(slug string, limit int, since string, desc bool, db *database.DB) Users {
	var (
		rows *sql.Rows
		err  error
	)
	/*
		subQuery :=
			`select u.* from forum_user u
			JOIN
			((select distinct author from thread WHERE LOWER(forum)=LOWER($1))
			UNION
			(select distinct author from post WHERE LOWER(forum)=LOWER($1))) as p ON nickname=p.author
			ORDER BY LOWER(nickname) `
		if desc {
			subQuery += "DESC "
		}
		if since != "" {
			query := "SELECT * FROM (" + subQuery + ") as sub WHERE LOWER(sub.nickname) "
			if desc {
				query += "< LOWER($2) "
			} else {
				query += "> LOWER($2) "
			}
			if limit != 0 {
				rows, err = db.DataBase.Query(query+"LIMIT $3;", slug, since, limit)
			} else {
				rows, err = db.DataBase.Query(query+";", slug, since)
			}
		} else {
			if limit != 0 {
				rows, err = db.DataBase.Query(subQuery+"LIMIT $2;", slug, limit)
			} else {
				rows, err = db.DataBase.Query(subQuery, slug)
			}
		}*/

	// Нет вложенного
	subQuery :=
		`select u.* from forum_user u
		JOIN 
		((select distinct author from thread WHERE LOWER(forum)=LOWER($1))
		UNION 
		(select distinct author from post WHERE LOWER(forum)=LOWER($1))) as p ON nickname=p.author `
	if since != "" {
		subQuery += "WHERE LOWER(u.nickname) "
		if desc {
			subQuery += "< LOWER($2) "
		} else {
			subQuery += "> LOWER($2) "
		}
		subQuery += "ORDER BY LOWER(nickname) "
		if desc {
			subQuery += "DESC "
		}
		if limit != 0 {
			rows, err = db.DataBase.Query(subQuery+"LIMIT $3;", slug, since, limit)
		} else {
			rows, err = db.DataBase.Query(subQuery+";", slug, since)
		}
	} else {
		subQuery += "ORDER BY LOWER(nickname) "
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
