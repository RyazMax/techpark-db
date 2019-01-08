package models

import (
	"database/sql"
	"log"
	"techpark-db/database"

	"github.com/astaxie/beego"
)

type Post struct {
	Author   string `json:"author"`
	Created  string `json:"created"`
	Forum    string `json:"forum"`
	Id       int    `json:"id"`
	IsEdited bool   `json:"isEdited"`
	Message  string `json:"message"`
	Parent   int    `json:"parent"`
	Thread   int    `json:"thread"`
	Mpath    []int  `json:"mpath,ommitempty"`
}

type PostFull struct {
	Author *User   `json:"author,ommitempty"`
	Forum  *Forum  `json:"forum,ommitempty"`
	Post   *Post   `json:"post"`
	Thread *Thread `json:"thread,ommitempty"`
}

type PostUpdate struct {
	Message string `json:"message,ommitempty"`
}

type Posts []Post

func (p *Post) Add(db *database.DB) error {
	//var d []uint8
	err := db.DataBase.QueryRow("insert into post(author,msg,parent,forum,thread,created)"+
		"values ($1,$2,$3,$4,$5,$6) RETURNING author,created,forum,id,isedited,msg,parent,thread;",
		p.Author, p.Message, p.Parent, p.Forum, p.Thread, p.Created).
		Scan(&p.Author, &p.Created, &p.Forum, &p.Id, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
	if err != nil {
		beego.Warn(err)
		log.Println(err)
	}
	return err
}

func (p *Post) Update(db *database.DB) {
	_, err := db.DataBase.Exec("UPDATE post SET msg=$1, isedited=true WHERE id=$2;", p.Message, p.Id)
	if err != nil {
		beego.Warn(err)
	}
}

func (p *Post) GetMpath(db *database.DB) (res []int) {
	err := db.DataBase.QueryRow("SELECT mpath FROM post WHERE id=$1;", p.Id).Scan(&res)
	if err != nil {
		beego.Warn(err)
	}
	return res
}

func (p *Post) GetByID(id int, db *database.DB) bool {
	err := db.DataBase.QueryRow("SELECT author,created,forum,id,isedited,msg,parent,thread FROM post WHERE id=$1;", id).
		Scan(&p.Author, &p.Created, &p.Forum, &p.Id, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
	if err != nil {
		return false
	}
	return true
}

func GetPostsSortedFlat(id int, limit int, since string, desc bool, db *database.DB) Posts {
	var (
		rows *sql.Rows
		err  error
	)

	/*
		subQuery := "SELECT * FROM post WHERE thread = $1 ORDER BY created, id "
		if desc {
			subQuery += "DESC "
		}
		if since != "" {
			query := "SELECT * FROM (" + subQuery + ") as sub WHERE sub.id "
			if desc {
				query += "< $2 "
			} else {
				query += "> $2 "
			}
			if limit != 0 {
				rows, err = db.DataBase.Query(query+"LIMIT $3;", id, since, limit)
			} else {
				rows, err = db.DataBase.Query(query+";", id, since)
			}
		} else {
			if limit != 0 {
				rows, err = db.DataBase.Query(subQuery+"LIMIT $2;", id, limit)
			} else {
				rows, err = db.DataBase.Query(subQuery, id)
			}
		}*/

	subQuery := "SELECT author,created,forum,id,isedited,msg,parent,thread FROM post WHERE thread = $1"
	if since != "" {
		subQuery += " AND id "
		if desc {
			subQuery += "< $2 "
		} else {
			subQuery += "> $2 "
		}
		subQuery += "ORDER BY id "
		if desc {
			subQuery += "DESC "
		}
		if limit != 0 {
			rows, err = db.DataBase.Query(subQuery+"LIMIT $3;", id, since, limit)
		} else {
			rows, err = db.DataBase.Query(subQuery+";", id, since)
		}
	} else {
		subQuery += "ORDER BY id "
		if desc {
			subQuery += "DESC "
		}
		if limit != 0 {
			rows, err = db.DataBase.Query(subQuery+"LIMIT $2;", id, limit)
		} else {
			rows, err = db.DataBase.Query(subQuery, id)
		}
	}
	defer rows.Close()
	if err != nil {
		beego.Warn(err)
		return nil
	}
	posts := make(Posts, 0)
	for rows.Next() {
		var t Post
		err = rows.Scan(&t.Author, &t.Created, &t.Forum, &t.Id, &t.IsEdited, &t.Message, &t.Parent, &t.Thread)
		if err != nil {
			beego.Warn(err)
			return nil
		}
		posts = append(posts, t)
	}
	return posts
}

func GetPostsSortedTree(id int, limit int, since string, desc bool, db *database.DB) Posts {
	var (
		rows *sql.Rows
		err  error
	)

	query := `
	SELECT author, created,forum,id,isedited,msg,parent,thread FROM post WHERE thread=$1 `
	if since != "" {
		query = "WITH since AS (SELECT mpath FROM post WHERE id=$2) " + query
		query += "AND mpath "
		if desc {
			query += "< (SELECT mpath FROM since) "
		} else {
			query += "> (SELECT mpath FROM since) "
		}
		query += "ORDER BY mpath "
		if desc {
			query += "DESC"
		}
		if limit != 0 {
			query += " LIMIT $3;"
			rows, err = db.DataBase.Query(query, id, since, limit)
		} else {
			query += ";"
			rows, err = db.DataBase.Query(query, id, since)
		}
	} else {
		query += "ORDER BY mpath "
		if desc {
			query += "DESC"
		}
		//query += ",mpath[2:] "
		if limit != 0 {
			query += " LIMIT $2;"
			rows, err = db.DataBase.Query(query, id, limit)
		} else {
			query += ";"
			rows, err = db.DataBase.Query(query, id)
		}
	}

	defer rows.Close()
	if err != nil {
		beego.Warn(err)
		return nil
	}
	posts := make(Posts, 0)
	for rows.Next() {
		var t Post
		err = rows.Scan(&t.Author, &t.Created, &t.Forum, &t.Id, &t.IsEdited, &t.Message, &t.Parent, &t.Thread)
		if err != nil {
			beego.Warn(err)
			return nil
		}
		posts = append(posts, t)
	}
	return posts
}

func GetPostsSortedParentTree(id int, limit int, since string, desc bool, db *database.DB) Posts {
	var (
		rows *sql.Rows
		err  error
	)
	// sorted AS (SELECT * FROM post WHERE thread=B366PapXAi86r , pag AS (SELECT mpath FROM sorted WHERE parent=0 OFFSET $2 LIMIT 1)
	query := `
	sorted AS (SELECT * FROM post WHERE thread=$1 `
	if since != "" {
		query = "since AS (SELECT mpath FROM post WHERE id=$2), " + query
		query += "AND mpath "
		if desc {
			query += "< (SELECT mpath FROM since) "
		} else {
			query += "> (SELECT mpath FROM since) "
		}
		query += "ORDER BY mpath "
		if desc {
			query += "DESC"
		}
		query += ")"
	} else {
		query += "ORDER BY mpath "
		if desc {
			query += "DESC"
		}
		query += ")"
	}

	query = "WITH " + query
	if limit != 0 {
		if since != "" {
			query += ", pag AS (SELECT mpath FROM sorted WHERE parent=0 OFFSET $3 LIMIT 1) "
		} else {
			query += ", pag AS (SELECT mpath FROM sorted WHERE parent=0 OFFSET $2 LIMIT 1) "
		}
		query += "SELECT author, created,forum,id,isedited,msg,parent,thread FROM sorted" // WHERE NOT mpath "
		/*if desc {
			query += " < "
		} else {
			query += " > "
		}
		*/
		//query += " (SELECT COALESCE((SELECT mpath FROM pag), ARRAY[]::bigint[])) OR mpath[0] = (SELECT mpath[0] FROM pag) OR (SELECT COALESCE((SELECT mpath FROM pag), ARRAY[]::bigint[])) = ARRAY[]::bigint[];"
		if since != "" {
			rows, err = db.DataBase.Query(query, id, since, limit)
		} else {
			rows, err = db.DataBase.Query(query, id, limit)
		}
	} else {
		query += "SELECT author, created,forum,id,isedited,msg,parent,thread FROM sorted;"
		if since != "" {
			rows, err = db.DataBase.Query(query, id, since)
		} else {
			rows, err = db.DataBase.Query(query, id)
		}
	}

	defer rows.Close()
	if err != nil {
		beego.Warn(err)
		return nil
	}
	posts := make(Posts, 0)
	for rows.Next() {
		var t Post
		err = rows.Scan(&t.Author, &t.Created, &t.Forum, &t.Id, &t.IsEdited, &t.Message, &t.Parent, &t.Thread)
		if err != nil {
			beego.Warn(err)
			return nil
		}
		posts = append(posts, t)
	}
	return posts
}
