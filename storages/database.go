package storages

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"sync-bot/config"
)

type Database struct {
	DB *sql.DB
}

func NewDBConnection(conf *config.Config) *Database {
	db, err := sql.Open("sqlite3", fmt.Sprintf("./%s", conf.Database.Path))
	if err != nil {
		log.Panicf("Error has occured during open database: %e", err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	database := &Database{db}
	return database
}
