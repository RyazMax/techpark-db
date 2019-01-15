package database

import (
	"io/ioutil"
	"log"

	"github.com/astaxie/beego"

	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type DB struct {
	DataBase *pgx.ConnPool
}

/*func (db *DB) ConectDB() {
	connStr := "user=docker password=docker dbname=docker sslmode=disable"

	newDB, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = newDB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	db.DataBase = newDB
	log.Println("Database conected")
}*/

func (db *DB) GetPool() {
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "docker",
			Password: "docker",
			Database: "docker",
		},
		MaxConnections: 60,
	})
	if err != nil {
		beego.Error(err)
	}
	db.DataBase = connPool
}

func (db DB) InitDB(filename string) {
	pd, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	cmd := string(pd)
	res, err := db.DataBase.Exec(cmd)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)
}
