package repository

import (
	"go.uber.org/zap"

	sq "github.com/Masterminds/squirrel"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3"

	"sync-bot/models"
	"sync-bot/storages"
)

type UserRepository interface {
	SaveNewUser(user models.User) error
	UpdateUserRank(user models.UpdatedUser) error
}

type userRepository struct {
	*storages.Database
}

func NewUserRepository(db *storages.Database) UserRepository {
	return &userRepository{
		db,
	}
}

func (r *userRepository) SaveNewUser(user models.User) error {
	query := sq.Insert("users").
		Columns("name", "rank").
		Values(user.Name, user.Rank)
	sql, params, err := query.ToSql()
	if err != nil {
		return err
	}
	result, err := r.DB.Exec(sql, params...)
	if sqliteErr, ok := err.(sqlite3.Error); ok {
		if sqliteErr.Code == sqlite3.ErrConstraint {
			return nil
		}
		return err
	}
	zap.S().Info("SaveNewUser: %v", result)
	return nil
}

func (r *userRepository) UpdateUserRank(user models.UpdatedUser) error {
	query := sq.Update("users").
		Set("rank", user.Rank).
		Where(sq.Eq{"name": user.Name})
	sql, params, err := query.ToSql()
	if err != nil {
		return err
	}
	result, err := r.DB.Exec(sql, params...)
	if err != nil {
		return err
	}
	zap.S().Info("UpdateUserRank: %v", result)
	return nil
}
