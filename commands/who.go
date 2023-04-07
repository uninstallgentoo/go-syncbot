package commands

import (
	"fmt"
	"math/rand"
	"strings"

	"github.com/uninstallgentoo/go-syncbot/command"
	"github.com/uninstallgentoo/go-syncbot/models"
)

var Who = &command.Command{
	Name:        "who",
	Description: "pick random user from userlist",
	Rank:        1,
	ExecFunc: func(args []string, c *command.Command) (models.CommandResult, error) {
		users := c.Processors.Chat.GetUsers()
		keys := make([]string, len(users))

		i := 0
		for k := range users {
			keys[i] = k
			i++
		}
		if strings.HasSuffix(args[len(args)-1], "?") {
			args[len(args)-1] = strings.TrimSuffix(args[len(args)-1], "?")
		}
		choice := fmt.Sprintf("%s %s", keys[rand.Intn(len(keys))], strings.Join(args, " "))
		return models.NewCommandResult(
			models.NewChatMessage(choice),
		), nil

	},
	ValidateFunc: func(args []string) error {
		if len(args) < 1 {
			return atLeastTwoArgsError
		}
		return nil
	},
}
