package client

import (
	"sync-bot/pkg/commands"
	"sync-bot/pkg/models"
	"sync-bot/pkg/service"
)

type Chat struct {
	chatService    service.Chat
	userService    service.UserService
	commandHandler commands.CommandHandler
	messages       []*models.Message
	users          map[string]*models.User
	commandResults chan commands.Event
}

func NewChatHandler(chatService service.Chat, userService service.UserService, handler commands.CommandHandler) *Chat {
	return &Chat{
		chatService:    chatService,
		userService:    userService,
		commandHandler: handler,
		commandResults: make(chan commands.Event),
		messages:       make([]*models.Message, 0),
		users:          make(map[string]*models.User, 0),
	}
}

func (c *Chat) handleMessage(msg *models.Message) {
	cleanedMessage := msg.Clean()
	result, isCommand := c.commandHandler.Parse(cleanedMessage.Text)
	if isCommand {
		result := c.commandHandler.Execute(result.Expr, result.Args, c.users[msg.Username].Rank)
		if result != nil {
			for _, response := range result.Results {
				if response != nil {
					c.commandResults <- *response
				}
			}
		}
	} else {
		c.messages = append(c.messages, cleanedMessage)
	}
}
