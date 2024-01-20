package storages

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"github.com/uninstallgentoo/go-syncbot/config"
)

type Database struct {
	DB *sql.DB
}

func NewDBConnection(conf *config.Config) *Database {
	db, err := sql.Open("sqlite3", conf.Database.Path)
	if err != nil {
		log.Panicf("Error has occured during open database: %e", err)
	}
	if err = db.Ping(); err != nil {
		log.Panic(err)
	}
	database := &Database{db}
	return database
}
