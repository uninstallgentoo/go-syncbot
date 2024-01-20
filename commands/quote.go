package commands

import (
	"github.com/uninstallgentoo/go-syncbot/command"
	"github.com/uninstallgentoo/go-syncbot/models"
)

var Quote = &command.Command{
	Name:        "q",
	Description: "pick random user's message from the chat history",
	Rank:        1,
	ExecFunc: func(args []string, c *command.Command) (models.CommandResult, error) {
		username := models.User{
			Name: args[0],
		}
		quote, err := c.Processors.Chat.GetUserQuote(username)
		if err != nil {
			return models.CommandResult{}, err
		}
		return models.NewCommandResult(
			models.NewChatMessage(quote),
		), nil

	},
	ValidateFunc: func(args []string) error {
		if len(args) < 1 {
			return atLeastTwoArgsError
		}
		return nil
	},
}
