package processors

import (
	"github.com/uninstallgentoo/go-syncbot/models"
	"github.com/uninstallgentoo/go-syncbot/repository"
)

type ChatHandler interface {
	SaveMessage(msg models.Message) error
	AddUserToList(users []models.User)
	SaveNewUser(user models.User) error
	DeleteUser(user models.UserLeave)
	UpdateUserRank(user models.UpdatedUser) error
	GetCommandResults() chan Event
	GetUsers() map[string]models.User
}

type Chat struct {
	chatRepo       repository.SyncRepository
	userRepo       repository.UserRepository
	users          map[string]models.User
	commandResults chan Event
}

func NewChatHandler(repositories *repository.Repositories) *Chat {
	return &Chat{
		chatRepo:       repositories.Sync,
		userRepo:       repositories.User,
		commandResults: make(chan Event),
		users:          make(map[string]models.User, 0),
	}
}

func (c *Chat) SaveMessage(msg models.Message) error {
	return c.chatRepo.SaveHistory(msg)
}

func (c *Chat) SaveNewUser(user models.User) error {
	if err := c.userRepo.SaveNewUser(user); err != nil {
		return err
	}
	if _, ok := c.users[user.Name]; !ok {
		c.users[user.Name] = user
	}
	return nil
}

func (c *Chat) AddUserToList(users []models.User) {
	for _, user := range users {
		c.users[user.Name] = user
	}
}

func (c *Chat) DeleteUser(user models.UserLeave) {
	if _, ok := c.users[user.Name]; ok {
		delete(c.users, user.Name)
	}
}

func (c *Chat) UpdateUserRank(user models.UpdatedUser) error {
	if err := c.userRepo.UpdateUserRank(user); err != nil {
		return err
	}
	return nil
}

func (c *Chat) GetCommandResults() chan Event {
	return c.commandResults
}

func (c *Chat) GetUsers() map[string]models.User {
	return c.users
}
