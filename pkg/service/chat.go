package service

import (
	"sync-bot/pkg/models"
	"sync-bot/pkg/repository"
)

type Chat interface {
	SaveChatHistory(messages []*models.Message) error
	GetRandomUserMessage(username string) (string, error)
}

type service struct {
	repo repository.SyncRepository
}

func NewService(repo repository.SyncRepository) Chat {
	return &service{repo}
}

func (s *service) SaveChatHistory(messages []*models.Message) error {
	return s.repo.SaveHistory(messages)
}

func (s *service) GetRandomUserMessage(username string) (string, error) {
	return s.repo.RandomMessage(username)
}
