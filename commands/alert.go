package commands

import (
	"strings"

	"github.com/uninstallgentoo/go-syncbot/command"
	"github.com/uninstallgentoo/go-syncbot/models"
)

var Alert = &command.Command{
	Name:        "alert",
	Description: "alert all users which in afk state",
	Rank:        2,
	ExecFunc: func(args []string, cmd *command.Command) (models.CommandResult, error) {
		usersToAlert := make([]string, 0)

		users := cmd.Processors.Chat.GetUsers()

		for _, user := range users {
			if user.Meta.AFK {
				usersToAlert = append(usersToAlert, user.Name)
			}
		}

		return models.NewCommandResult(
			models.NewChatMessage(strings.Join(usersToAlert, " ")),
		), nil
	},
}
