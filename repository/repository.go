package repository

import (
	"fmt"
	"github.com/uninstallgentoo/go-syncbot/models"

	sq "github.com/Masterminds/squirrel"

	"github.com/uninstallgentoo/go-syncbot/storages"
)

type SyncRepository interface {
	SaveHistory(message models.Message) error
	RandomMessage(username string) (string, error)
	FetchUserStatistic(username string) (int, error)
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

func (r *syncRepository) FetchUserStatistic(username string) (int, error) {
	query := sq.Select("count(*)").
		From("chat_history").
		Where(sq.Eq{"username": username})
	sql, params, err := query.ToSql()
	if err != nil {
		return 0, err
	}
	var count int
	err = r.DB.QueryRow(sql, params...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
