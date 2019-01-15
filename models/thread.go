package models

import (
	"errors"
	"fmt"
	"strings"
	"techpark-db/database"
	"time"

	"github.com/jackc/pgx"

	"github.com/astaxie/beego"
)

type Thread struct {
	Author   string    `json:"author"`
	Created  time.Time `json:"created,ommitempty"`
	Forum    string    `json:"forum"`
	ID       int       `json:"id"`
	IsEdited bool      `json:"isEdited"`
	Message  string    `json:"message"`
	Slug     string    `json:"slug,ommitempty"`
	Title    string    `json:"title"`
	Votes    int       `json:"votes"`
}

type ThreadUpdate struct {
	Message string `json:"message"`
	Title   string `json:"title"`
}

type Threads []Thread

func (t *Thread) Add(db *database.DB) {
	var err error
	if t.Created.IsZero() {
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

func (t *Thread) GetVotesById(id int, db *database.DB) (votes int) {
	err := db.DataBase.QueryRow("SELECT votes FROM THREAD WHERE id = $1;", id).
		Scan(&votes)
	if err != nil {
		beego.Warn(err)
		return 0
	}
	return votes
}

func (t *Thread) GetBySlug(slug string, db *database.DB) bool {
	err := db.DataBase.QueryRow("SELECT * FROM THREAD WHERE slug=$1", slug).
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

func (t *Thread) AddPosts(posts Posts, db *database.DB) ([]int, time.Time, error) {
	result := make([]int, 0)
	curTime := time.Now()
	//thread_ids := t.GetPostsID(db)
	authors := make(map[string]bool)
	parents := make(map[int]bool)
	for _, post := range posts {
		post.Thread = t.ID
		post.Forum = t.Forum
		authors[post.Author] = true
		if post.Parent != 0 {
			parents[post.Parent] = true
		}
	}

	tmp := GetUsersByNicks(&authors, db)
	if len(tmp) != len(authors) {
		return result, curTime, errors.New("No author")
	}
	parents_found := GetPostsByID(&parents, t.ID, db)
	if parents_found != len(parents) {
		return result, curTime, errors.New("Parent in other thread")
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
	query.WriteString("RETURNING id;")
	if len(posts) > 0 {
		rows, _ := db.DataBase.Query(query.String(), args...)
		defer rows.Close()
		/*if err != nil {
			pqErr := err.(*pq.Error)
			beego.Warn(pqErr.InternalQuery)
			beego.Warn(pqErr.Error())
			return posts, err
		}*/
		var id int
		for rows.Next() {

			rows.Scan(&id)
			result = append(result, id)
		}
	}

	//beego.Info("In model ", len(result))
	return result, curTime, nil
}

func GetThreadsSorted(slug string, limit int, since string, desc bool, db *database.DB) Threads {
	var (
		rows *pgx.Rows
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
	var subQuery strings.Builder
	subQuery.WriteString("SELECT * FROM THREAD WHERE FORUM = $1")
	if since != "" {
		subQuery.WriteString(" AND created ")
		if desc {
			subQuery.WriteString("<= $2 ")
		} else {
			subQuery.WriteString(">= $2 ")
		}
		subQuery.WriteString("ORDER BY created ")
		if desc {
			subQuery.WriteString("DESC ")
		}
		if limit != 0 {
			subQuery.WriteString("LIMIT $3;")
			rows, err = db.DataBase.Query(subQuery.String(), slug, since, limit)
		} else {
			subQuery.WriteString(";")
			rows, err = db.DataBase.Query(subQuery.String(), slug, since)
		}
	} else {
		subQuery.WriteString("ORDER BY created ")
		if desc {
			subQuery.WriteString("DESC ")
		}
		if limit != 0 {
			subQuery.WriteString("LIMIT $2;")
			rows, err = db.DataBase.Query(subQuery.String(), slug, limit)
		} else {
			rows, err = db.DataBase.Query(subQuery.String(), slug)
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
