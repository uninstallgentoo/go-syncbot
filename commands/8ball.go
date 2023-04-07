package commands

import (
	"math/rand"

	"github.com/uninstallgentoo/go-syncbot/command"
	"github.com/uninstallgentoo/go-syncbot/models"
)

var choices = []string{
	"It is certain.",
	"It is decidedly so.",
	"Without a doubt.",
	"Yes definitely.",
	"You may rely on it.",
	"As I see it, yes.",
	"Most likely.",
	"Outlook good.",
	"Yes.",
	"Signs point to yes.",
	"Reply hazy, try again.",
	"Ask again later.",
	"Better not tell you now.",
	"Cannot predict now.",
	"Concentrate and ask again.",
	"Don't count on it.",
	"My reply is no.",
	"My sources say no.",
	"Outlook not so good.",
	"Very doubtful.",
}

var MagicBall = &command.Command{
	Name:        "8ball",
	Description: "Get random response to your answer.",
	Rank:        1,
	ExecFunc: func(args []string, c *command.Command) (models.CommandResult, error) {
		choice := choices[rand.Intn(len(choices))]
		return models.NewCommandResult(
			models.NewChatMessage(choice),
		), nil
	},
}
