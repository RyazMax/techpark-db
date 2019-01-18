package models

import (
	"fmt"
	"log"
	"strings"
	"techpark-db/database"
	"time"

	"github.com/astaxie/beego"
	"github.com/jackc/pgx"
)

type Post struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created"`
	Forum    string    `json:"forum"`
	Id       int       `json:"id"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Parent   int       `json:"parent"`
	Thread   int       `json:"thread"`
	Mpath    []int     `json:"mpath,ommitempty"`
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

//easyjson:json
type Posts []Post

func PostUpd(p Post, db *database.DB) {
	_, err := db.DataBase.Exec("UPDATE post SET msg=$1, isedited=true WHERE id=$2;", p.Message, p.Id)
	if err != nil {
		log.Println(err)
	}
}

func (p *Post) GetMpath(db *database.DB) (res []int) {
	err := db.DataBase.QueryRow("SELECT mpath FROM post WHERE id=$1;", p.Id).Scan(&res)
	if err != nil {
		beego.Warn(err)
	}
	return res
}

func PostGetByID(id int, db *database.DB) (p Post, exist bool) {
	err := db.DataBase.QueryRow("SELECT author,created,forum,id,isedited,msg,parent,thread FROM post WHERE id=$1;", id).
		Scan(&p.Author, &p.Created, &p.Forum, &p.Id, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
	if err != nil {
		return Post{}, false
	}
	return p, true
}

func GetPostsByID(ids *map[int]bool, thread int, db *database.DB) (res int) {
	if len(*ids) == 0 {
		return 0
	}
	var query strings.Builder
	args := make([]interface{}, 0, len(*ids))
	query.Grow(70 + len(*ids)*4)
	query.WriteString("SELECT id FROM post WHERE id in (")
	var cnt int
	for id := range *ids {
		if cnt > 0 {
			query.WriteString(", ")
		}
		cnt++
		query.WriteString(fmt.Sprintf("$%d", cnt))
		args = append(args, id)
	}
	query.WriteString(fmt.Sprintf(") AND thread = $%d;", len(*ids)+1))
	args = append(args, thread)

	rows, err := db.DataBase.Query(query.String(), args...)
	defer rows.Close()
	if err != nil {
		beego.Warn(err)
		return 0
	}
	var tmp int
	for rows.Next() {
		rows.Scan(&tmp)
		res++
	}
	return
}

func GetPostsSortedFlat(id int, limit int, since string, desc bool, db *database.DB) Posts {
	var (
		rows *pgx.Rows
		err  error
	)
	var subQuery strings.Builder
	subQuery.WriteString("SELECT author,created,forum,id,isedited,msg,parent,thread FROM post WHERE thread = $1")
	if since != "" {
		subQuery.WriteString(" AND id ")
		if desc {
			subQuery.WriteString("< $2 ")
		} else {
			subQuery.WriteString("> $2 ")
		}
		subQuery.WriteString("ORDER BY id ")
		if desc {
			subQuery.WriteString("DESC ")
		}
		if limit != 0 {
			subQuery.WriteString("LIMIT $3;")
			rows, err = db.DataBase.Query(subQuery.String(), id, since, limit)
		} else {
			subQuery.WriteString(";")
			rows, err = db.DataBase.Query(subQuery.String(), id, since)
		}
	} else {
		subQuery.WriteString("ORDER BY id ")
		if desc {
			subQuery.WriteString("DESC ")
		}
		if limit != 0 {
			subQuery.WriteString("LIMIT $2;")
			rows, err = db.DataBase.Query(subQuery.String(), id, limit)
		} else {
			subQuery.WriteString(";")
			rows, err = db.DataBase.Query(subQuery.String(), id)
		}
	}
	defer rows.Close()
	posts := make(Posts, 0)
	for rows.Next() {
		var t Post
		err = rows.Scan(&t.Author, &t.Created, &t.Forum, &t.Id, &t.IsEdited, &t.Message, &t.Parent, &t.Thread)
		if err != nil {
			return nil
		}
		posts = append(posts, t)
	}
	return posts
}

func GetPostsSortedTree(id int, limit int, since string, desc bool, db *database.DB) Posts {
	var (
		rows *pgx.Rows
		err  error
	)

	var query strings.Builder
	if since != "" {
		query.WriteString("WITH since AS (SELECT mpath FROM post WHERE id=$2) ")
	}
	query.WriteString(`
	SELECT author, created,forum,id,isedited,msg,parent,thread FROM post WHERE thread=$1 `)
	if since != "" {
		query.WriteString("AND mpath ")
		if desc {
			query.WriteString("< (SELECT mpath FROM since) ")
		} else {
			query.WriteString("> (SELECT mpath FROM since) ")
		}
		query.WriteString("ORDER BY mpath ")
		if desc {
			query.WriteString("DESC")
		}
		if limit != 0 {
			query.WriteString(" LIMIT $3;")
			rows, err = db.DataBase.Query(query.String(), id, since, limit)
		} else {
			query.WriteString(";")
			rows, err = db.DataBase.Query(query.String(), id, since)
		}
	} else {
		query.WriteString("ORDER BY mpath ")
		if desc {
			query.WriteString("DESC")
		}
		if limit != 0 {
			query.WriteString(" LIMIT $2;")
			rows, err = db.DataBase.Query(query.String(), id, limit)
		} else {
			query.WriteString(";")
			rows, err = db.DataBase.Query(query.String(), id)
		}
	}

	defer rows.Close()
	posts := make(Posts, 0)
	for rows.Next() {
		var t Post
		err = rows.Scan(&t.Author, &t.Created, &t.Forum, &t.Id, &t.IsEdited, &t.Message, &t.Parent, &t.Thread)
		if err != nil {
			return nil
		}
		posts = append(posts, t)
	}
	return posts
}

func GetPostsSortedParentTree(id int, limit int, since string, desc bool, db *database.DB) Posts {
	var (
		rows  *pgx.Rows
		err   error
		query strings.Builder
	)
	query.WriteString("WITH ")
	if since != "" {
		query.WriteString("since AS (SELECT mpath FROM post WHERE id=$2), ")
	}
	query.WriteString(`
	sorted AS (SELECT * FROM post WHERE thread=$1 `)
	if since != "" {
		query.WriteString("AND mpath[1] ")
		if desc {
			query.WriteString("< (SELECT mpath[1] FROM since) ")
		} else {
			query.WriteString("> (SELECT mpath[1] FROM since) ")
		}
		query.WriteString("ORDER BY mpath[1] ")
		if desc {
			query.WriteString("DESC")
		}
		query.WriteString(", mpath[1:]")
		query.WriteString(")")
	} else {
		query.WriteString("ORDER BY mpath[1] ")
		if desc {
			query.WriteString("DESC")
		}
		query.WriteString(", mpath[1:]")
		query.WriteString(")")
	}

	if limit != 0 {
		if since != "" {
			query.WriteString(", pag AS (SELECT mpath FROM sorted WHERE parent=0 OFFSET $3 LIMIT 1) ")
		} else {
			query.WriteString(", pag AS (SELECT mpath FROM sorted WHERE parent=0 OFFSET $2 LIMIT 1) ")
		}
		query.WriteString("SELECT author, created,forum,id,isedited,msg,parent,thread FROM sorted  WHERE NOT mpath ")
		if desc {
			query.WriteString(" < ")
		} else {
			query.WriteString(" > ")
		}

		query.WriteString(" (SELECT COALESCE((SELECT mpath FROM pag), ARRAY[]::integer[])) OR mpath[1] = (SELECT mpath[1] FROM pag) OR (SELECT COALESCE((SELECT mpath FROM pag), ARRAY[]::integer[])) = ARRAY[]::integer[];")
		if since != "" {
			rows, err = db.DataBase.Query(query.String(), id, since, limit-1)
		} else {
			rows, err = db.DataBase.Query(query.String(), id, limit-1)
		}
	} else {
		query.WriteString("SELECT author, created,forum,id,isedited,msg,parent,thread FROM sorted;")
		if since != "" {
			rows, err = db.DataBase.Query(query.String(), id, since)
		} else {
			rows, err = db.DataBase.Query(query.String(), id)
		}
	}

	defer rows.Close()
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
