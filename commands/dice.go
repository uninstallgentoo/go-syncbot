package commands

import (
	"errors"
	"math/rand"
	"strconv"

	"github.com/uninstallgentoo/go-syncbot/command"
	"github.com/uninstallgentoo/go-syncbot/models"
)

var diceInputIsNotANumberError = errors.New("input value is not a number")

var diceInputLessThanZeroError = errors.New("input value is less than 0")

var Dice = &command.Command{
	Name:        "dice",
	Description: "roll a dice",
	Rank:        1,
	ExecFunc: func(args []string, cmd *command.Command) (models.CommandResult, error) {
		payloads := make([]*models.Event, 0)

		val, _ := strconv.Atoi(args[0])

		payloads = append(payloads, &models.Event{
			Method:  "chatMsg",
			Message: models.EventPayload{Message: strconv.Itoa(rand.Intn(val)), Meta: struct{}{}},
		})
		return models.CommandResult{Results: payloads}, nil
	},
	ValidateFunc: func(args []string) error {
		var val int
		var err error
		val, err = strconv.Atoi(args[0])
		if err != nil {
			return diceInputIsNotANumberError
		}
		if val <= 0 {
			return diceInputLessThanZeroError
		}
		return nil
	},
}
