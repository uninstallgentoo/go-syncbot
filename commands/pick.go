package commands

import (
	"errors"
	"math/rand"

	"github.com/uninstallgentoo/go-syncbot/command"
	"github.com/uninstallgentoo/go-syncbot/models"
)

var atLeastTwoArgsError = errors.New("There should be at least to values to pick from")

var Pick = &command.Command{
	Name:        "pick",
	Description: "pick random value from the list of arguments",
	Rank:        1,
	ExecFunc: func(args []string, c *command.Command) (models.CommandResult, error) {
		choice := args[rand.Intn(len(args))]
		return models.NewCommandResult(
			models.NewChatMessage(choice),
		), nil
	},
	ValidateFunc: func(args []string) error {
		if len(args) < 2 {
			return atLeastTwoArgsError
		}
		return nil
	},
}
