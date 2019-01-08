package models

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"techpark-db/database"
	"time"

	"github.com/astaxie/beego"
)

type Thread struct {
	Author   string `json:"author"`
	Created  string `json:"created,ommitempty"`
	Forum    string `json:"forum"`
	ID       int    `json:"id"`
	IsEdited bool   `json:"isEdited"`
	Message  string `json:"message"`
	Slug     string `json:"slug,ommitempty"`
	Title    string `json:"title"`
	Votes    int    `json:"votes"`
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type Threads []Thread

func (t *Thread) Add(db *database.DB) {
	var err error
	if t.Created == "" {
		err = db.DataBase.QueryRow("INSERT INTO THREAD(author, forum, msg, title, slug) values ($1, $2, $3, $4, $5) RETURNING *;", t.Author, t.Forum, t.Message, t.Title, t.Slug).
			Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.IsEdited, &t.Message, &t.Slug, &t.Title, &t.Votes)
	} else {
		err = db.DataBase.QueryRow("INSERT INTO THREAD(author, created, forum, msg, title, slug) values ($1, $2, $3, $4, $5, $6) RETURNING *;", t.Author, t.Created, t.Forum, t.Message, t.Title, t.Slug).
			Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.IsEdited, &t.Message, &t.Slug, &t.Title, &t.Votes)
	}
	if err != nil {
		beego.Warn(err)
		return
	}
}

func (t *Thread) Update(db *database.DB) {
	_, err := db.DataBase.Exec("UPDATE thread SET votes=$1, title=$3, msg=$4 WHERE id=$2;", t.Votes, t.ID, t.Title, t.Message)
	if err != nil {
		beego.Warn(err)
	}
}

func (t *Thread) GetById(id int, db *database.DB) bool {
	err := db.DataBase.QueryRow("SELECT author, created, forum, id, isedited, msg, slug, title, votes FROM THREAD WHERE id = $1", id).
		Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.IsEdited, &t.Message, &t.Slug, &t.Title, &t.Votes)
	if err != nil {
		return false
	}
	return true
}

func (t *Thread) GetBySlug(slug string, db *database.DB) bool {
	err := db.DataBase.QueryRow("SELECT * FROM THREAD WHERE LOWER(slug)=LOWER($1)", slug).
		Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.IsEdited, &t.Message, &t.Slug, &t.Title, &t.Votes)
	if err != nil {
		return false
	}
	return true
}

func (t *Thread) GetPostsID(db *database.DB) (res []int) {
	rows, err := db.DataBase.Query("SELECT id FROM post WHERE thread=$1;", t.ID)
	defer rows.Close()
	if err != nil {
		beego.Warn(err)
		return make([]int, 0)
	}
	var i int
	for rows.Next() {
		rows.Scan(&i)
		if err != nil {
			beego.Warn(err)
			return make([]int, 0)
		}
		res = append(res, i)
	}
	return
}

func (t *Thread) AddPosts(posts Posts, db *database.DB) (Posts, error) {
	result := make(Posts, 0)
	curTime := time.Now().Format(time.RFC3339)
	thread_ids := t.GetPostsID(db)
	for _, post := range posts {
		post.Thread = t.ID
		post.Forum = t.Forum
		author := User{}
		exist := author.GetUserByNick(post.Author, db)
		if !exist {
			return posts, errors.New("No author")
		}
		if post.Parent == 0 {
			continue
		}

		exist = false
		if post.Parent != 0 {
			for _, id := range thread_ids {
				if id == post.Parent {
					exist = true
				}
			}
		} else {
			exist = true
		}
		if !exist {
			return result, errors.New("No parent in thread")
		}
	}

	var query strings.Builder
	args := make([]interface{}, 0)
	query.WriteString("insert into post(author,msg,parent,forum,thread,created) values ")
	for i, post := range posts {
		post.Thread = t.ID
		post.Forum = t.Forum
		post.Created = curTime
		if i != 0 {
			query.WriteString(fmt.Sprintf(",($%d, $%d, $%d, $%d, $%d, $%d) ", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		} else {
			query.WriteString(fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d) ", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		}
		args = append(args, post.Author, post.Message, post.Parent, post.Forum, post.Thread, post.Created)
	}
	query.WriteString("RETURNING author,created,forum,id,isedited,msg,parent,thread;")
	if len(posts) > 0 {
		rows, err := db.DataBase.Query(query.String(), args...)
		defer rows.Close()
		if err != nil {
			beego.Warn(err)
			return posts, err
		}
		p := Post{}
		for rows.Next() {

			rows.Scan(&p.Author, &p.Created, &p.Forum, &p.Id, &p.IsEdited, &p.Message, &p.Parent, &p.Thread)
			result = append(result, p)
		}
	}

	return result, nil
}

func GetThreadsSorted(slug string, limit int, since string, desc bool, db *database.DB) Threads {
	var (
		rows *sql.Rows
		err  error
	)
	/*subQuery := "SELECT * FROM THREAD WHERE LOWER(FORUM) = LOWER($1) ORDER BY created "
	if desc {
		subQuery += "DESC "
	}
	if since != "" {
		query := "SELECT * FROM (" + subQuery + ") as sub WHERE sub.created "
		if desc {
			query += "<= $2 "
		} else {
			query += ">= $2 "
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

	// Исправлена вложенность
	subQuery := "SELECT * FROM THREAD WHERE LOWER(FORUM) = LOWER($1)"
	if since != "" {
		subQuery += " AND created "
		if desc {
			subQuery += "<= $2 "
		} else {
			subQuery += ">= $2 "
		}
		subQuery += "ORDER BY created "
		if desc {
			subQuery += "DESC "
		}
		if limit != 0 {
			rows, err = db.DataBase.Query(subQuery+"LIMIT $3;", slug, since, limit)
		} else {
			rows, err = db.DataBase.Query(subQuery+";", slug, since)
		}
	} else {
		subQuery += "ORDER BY created "
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
	threads := make(Threads, 0)
	for rows.Next() {
		var t Thread
		err = rows.Scan(&t.Author, &t.Created, &t.Forum, &t.ID, &t.IsEdited, &t.Message, &t.Slug, &t.Title, &t.Votes)
		if err != nil {
			beego.Warn(err)
			return nil
		}
		threads = append(threads, t)
	}
	return threads
}
