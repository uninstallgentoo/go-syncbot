package repository

import (
	"fmt"
	"sync-bot/models"

	sq "github.com/Masterminds/squirrel"

	"sync-bot/storages"
)

type SyncRepository interface {
	SaveHistory(message models.Message) error
	RandomMessage(username string) (string, error)
}

type syncRepository struct {
	*storages.Database
}

func NewSyncRepository(db *storages.Database) SyncRepository {
	return &syncRepository{
		db,
	}
}

func (r *syncRepository) SaveHistory(message models.Message) error {
	query := sq.Insert("chat_history").
		Columns("timestamp", "username", "msg").
		Values(message.Time, message.Username, message.Text)
	sql, params, err := query.ToSql()
	if err != nil {
		return err
	}
	if _, err = r.DB.Exec(sql, params...); err != nil {
		return err
	}
	return nil
}

func (r *syncRepository) RandomMessage(username string) (string, error) {
	var err error
	query := sq.Select("msg", "timestamp").
		From("chat_history").
		Where(sq.Eq{"username": username}).
		OrderBy("RANDOM()").
		Limit(1)
	sql, params, err := query.ToSql()
	if err != nil {
		return "", err
	}
	var timestamp, msg string
	err = r.DB.QueryRow(sql, params...).Scan(&msg, &timestamp)
	if err != nil {
		return "", err
	}
	quote := fmt.Sprintf("[%s][%s] %s", timestamp, username, msg)
	return quote, nil
}
