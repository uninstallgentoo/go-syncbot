package commands

import (
	"errors"
	"strconv"

	"github.com/uninstallgentoo/go-syncbot/command"
	"github.com/uninstallgentoo/go-syncbot/models"
)

var emptyUsernameError = errors.New("Username argument is required")
var usernameNotFoundError = errors.New("Username not found")

var Stat = &command.Command{
	Name:        "stat",
	Description: "fetch user chat statistic in channel",
	Rank:        1,
	ExecFunc: func(args []string, c *command.Command) (models.CommandResult, error) {
		username := models.User{
			Name: args[0],
		}
		count, err := c.Processors.Chat.FetchUserStatistic(username)
		if err != nil {
			return models.CommandResult{}, usernameNotFoundError
		}
		return models.NewCommandResult(
			models.NewChatMessage(strconv.Itoa(count)),
		), nil
	},
	ValidateFunc: func(args []string) error {
		if len(args) < 1 {
			return emptyUsernameError
		}
		return nil
	},
}
