package database

import (
	"database/sql"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

type DB struct {
	DataBase *sql.DB
}

func (db *DB) ConectDB() {
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
