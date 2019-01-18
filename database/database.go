package database

import (
	"io/ioutil"
	"log"

	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type DB struct {
	DataBase *pgx.ConnPool
}

var db DB

func (db *DB) GetPool() {
	connPool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Port:     5432,
			User:     "docker",
			Password: "docker",
			Database: "docker",
		},
		MaxConnections: 50,
	})
	if err != nil {
		log.Fatal(err)
	}
	db.DataBase = connPool
}

func (db DB) InitDB(filename string) {
	pd, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	cmd := string(pd)
	_, err = db.DataBase.Exec(cmd)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Database inited")
}

func GetDB() *DB {
	return &db
}
