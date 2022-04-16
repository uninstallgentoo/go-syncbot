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
		payloads := make([]*models.Event, 0)

		usersToAlert := make([]string, 0)

		users := cmd.Processors.Chat.GetUsers()
		for _, user := range users {
			if user.Meta.AFK {
				usersToAlert = append(usersToAlert, user.Name)
			}
		}

		payloads = append(payloads, &models.Event{
			Method:  "chatMsg",
			Message: models.EventPayload{Message: strings.Join(usersToAlert, " "), Meta: struct{}{}},
		})
		return models.CommandResult{Results: payloads}, nil
	},
}
