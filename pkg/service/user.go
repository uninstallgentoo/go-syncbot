package service

import (
	"sync-bot/pkg/models"
	"sync-bot/pkg/repository"
)

type UserService interface {
	SaveNewUser(user *models.User) error
	UpdateUserRank(*models.UpdatedUser) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) SaveNewUser(user *models.User) error {
	return s.repo.SaveNewUser(user)
}

func (s *userService) UpdateUserRank(user *models.UpdatedUser) error {
	return s.repo.UpdateUserRank(user)
}
