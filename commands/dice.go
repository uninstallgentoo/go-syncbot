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
var diceEmptyArgumentError = errors.New("you should provide value to roll a dice from")

var Dice = &command.Command{
	Name:        "dice",
	Description: "roll a dice",
	Rank:        1,
	ExecFunc: func(args []string, cmd *command.Command) (models.CommandResult, error) {
		val, _ := strconv.Atoi(args[0])
		return models.NewCommandResult(
			models.NewChatMessage(strconv.Itoa(rand.Intn(val))),
		), nil
	},
	ValidateFunc: func(args []string) error {
		var val int
		var err error
		if len(args) < 1 {
			return diceEmptyArgumentError
		}
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
