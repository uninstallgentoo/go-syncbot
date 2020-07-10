package repository

import (
	"fmt"
	"sync-bot/pkg/models"

	sq "github.com/Masterminds/squirrel"

	"sync-bot/pkg/storages"
)

type SyncRepository interface {
	SaveHistory([]*models.Message) error
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

func (r *syncRepository) SaveHistory(messages []*models.Message) error {
	query := sq.Insert("chat_history").
		Columns("timestamp", "username", "msg")
	for _, key := range messages {
		query = query.Values(key.Time, key.Username, key.Text)
	}
	sql, params, err := query.ToSql()
	if err != nil {
		return err
	}
	result, err := r.DB.Exec(sql, params...)
	if err != nil {
		return err
	}
	fmt.Printf("messages has been saved: %v", result)
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
