package storages

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"

	"sync-bot/pkg/config"
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
	database.initDatabase()
	return database
}

func (d *Database) initDatabase() {
	version, err := d.fetchLastVersion()
	if err != nil {
		d.initVersionTable()
	}
	if version != "1" {
		d.createUserTable()
		d.createChatHistoryTable()
	}
}

func (d *Database) initVersionTable() {
	d.createVersionTable()
	d.updateDBVersion(1)
}

func (d *Database) createVersionTable() {
	query := "CREATE TABLE IF NOT EXISTS version(value INTEGER)"
	_, err := d.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Database) updateDBVersion(version int) {
	query := "INSERT INTO version VALUES (?)"
	_, err := d.DB.Exec(query, version)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Database) createUserTable() {
	query := "CREATE TABLE IF NOT EXISTS users(name TEXT, rank INTEGER," +
		" PRIMARY KEY(name))"
	_, err := d.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Database) createChatHistoryTable() {
	query := "CREATE TABLE IF NOT EXISTS chat_history(timestamp INTEGER," +
		" username TEXT, msg TEXT)"
	_, err := d.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func (d *Database) fetchLastVersion() (string, error) {
	var value string
	err := d.DB.QueryRow("SELECT value FROM version ORDER BY value ASC LIMIT 1").Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil
}
