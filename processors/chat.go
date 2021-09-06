package processors

import (
	"sync"

	"github.com/uninstallgentoo/go-syncbot/models"
	"github.com/uninstallgentoo/go-syncbot/repository"
)

type ChatHandler interface {
	SaveMessage(msg models.Message) error
	AddUserToList(users []models.User)
	SaveNewUser(user models.User) error
	DeleteUser(user models.UserLeave)
	UpdateUserRank(user models.UpdatedUser) error
	GetUsers() map[string]models.User
}

type Chat struct {
	m        *sync.RWMutex
	chatRepo repository.SyncRepository
	userRepo repository.UserRepository
	users    map[string]models.User
}

func NewChatHandler(repositories *repository.Repositories) *Chat {
	return &Chat{
		m:        &sync.RWMutex{},
		chatRepo: repositories.Sync,
		userRepo: repositories.User,
		users:    make(map[string]models.User, 0),
	}
}

func (c *Chat) SaveMessage(msg models.Message) error {
	return c.chatRepo.SaveHistory(msg)
}

func (c *Chat) SaveNewUser(user models.User) error {
	if err := c.userRepo.SaveNewUser(user); err != nil {
		return err
	}
	c.m.RLock()
	defer c.m.RUnlock()
	if _, ok := c.users[user.Name]; !ok {
		c.users[user.Name] = user
	}
	return nil
}

func (c *Chat) AddUserToList(users []models.User) {
	c.m.RLock()
	defer c.m.RUnlock()
	for _, user := range users {
		c.users[user.Name] = user
	}
}

func (c *Chat) DeleteUser(user models.UserLeave) {
	c.m.RLock()
	defer c.m.RUnlock()
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

func (c *Chat) GetUsers() map[string]models.User {
	return c.users
}
