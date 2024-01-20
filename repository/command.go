package repository

import (
	"errors"
	"go.uber.org/zap"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"github.com/uninstallgentoo/go-syncbot/models"
	"github.com/uninstallgentoo/go-syncbot/storages"
)

type CommandRepository interface {
	UpdateRank(user *models.Command) error
	Add(user *models.Command) error
	FetchAll() ([]*models.Command, error)
}

type commandRepository struct {
	*storages.Database
}

func NewCommandRepository(db *storages.Database) CommandRepository {
	return &commandRepository{
		db,
	}
}

func (r *commandRepository) Add(user *models.Command) error {
	query := sq.Insert("command_rank").
		Columns("name", "rank").
		Values(user.Command, user.Rank)
	sql, params, err := query.ToSql()
	if err != nil {
		return err
	}
	result, err := r.DB.Exec(sql, params...)
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			return nil
		}
		return err
	}
	zap.S().Info("AddCommand: %v", result)
	return nil
}

func (r *commandRepository) UpdateRank(user *models.Command) error {
	query := sq.Update("command_rank").
		Set("rank", user.Rank).
		Where(sq.Eq{"command": user.Command})
	sql, params, err := query.ToSql()
	if err != nil {
		return err
	}
	result, err := r.DB.Exec(sql, params...)
	if err != nil {
		return err
	}
	zap.S().Info("UpdateCommandRank: %v", result)
	return nil
}

func (r *commandRepository) FetchAll() ([]*models.Command, error) {
	query := sq.Select("command", "rank").
		From("command_rank")
	sql, params, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := r.DB.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	commands := make([]*models.Command, 0)
	for rows.Next() {
		var (
			command string
			rank    float64
		)
		err := rows.Scan(&command, &rank)
		if err != nil {
			log.Fatal(err)
		}
		commands = append(commands, &models.Command{
			Command: command,
			Rank:    rank,
		})
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return commands, nil
}
