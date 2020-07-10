package commands

import (
	"errors"

	"sync-bot/pkg/service"
)

type randomMessageCommand struct {
	h service.Chat
}

func NewRandomMessageCommand(chatService service.Chat) CommandExecutor {
	return &randomMessageCommand{
		chatService,
	}
}

func (c *randomMessageCommand) GetMinRequiredRank() float64 {
	return 1
}

func (c *randomMessageCommand) Validate(args []string) error {
	if len(args) == 0 || args[0] == "" {
		return errors.New("Укажите ник пользователя.")
	}
	return nil
}

func (c *randomMessageCommand) Exec(args []string) (*CommandResult, error) {
	username := args[0]
	msg, err := c.h.GetRandomUserMessage(username)
	if err != nil {
		return nil, err
	}
	results := []*Event{
		{
			Method:  "chatMsg",
			Message: EventPayload{Message: msg, Meta: struct{}{}},
		},
	}
	return NewCommandResult(results), nil
}
