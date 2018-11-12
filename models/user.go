package models

import (
	"log"
	"techpark-db/database"
)

type User struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Nickname string `json:"nickname,omitempty"`
}

func (newUser *User) Add(db *database.DB) error {
	_, err := db.DataBase.Exec("insert into forum_user(email,about,fullname,nickname) values ($1,$2,$3,LOWER($4))		;",
		newUser.Email, newUser.About, newUser.Fullname, newUser.Nickname)
	if err != nil {

		return err
	}
	return nil
}

func (u *User) GetLike(db *database.DB) Users {
	rows, err := db.DataBase.Query("select * from forum_user where email=$1 or nickname=LOWER($2);",
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
	rows, err := db.DataBase.Query("select * from forum_user where nickname=LOWER($1);", nickname)
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
	_, err := db.DataBase.Exec("UPDATE forum_user SET(fullname=$1,about=$2,email=$3);", u.Fullname, u.About, u.Email)
	if err != nil {
		log.Println("Update", err)
	}
	return err
}

type UserUpdate struct {
	About    string `json:"about"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

type Users []User

//func (m *User) MarshalText() ([]byte, error) { return }

//func (m *User) UnmarshalText(b []byte) error { return }
