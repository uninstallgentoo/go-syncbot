package commands

import (
	"math/rand"
	"strconv"
)

type randCommand struct {
	AllowedRank float64
}

func NewRandCommand() CommandExecutor {
	return &randCommand{
		AllowedRank: 1,
	}
}

func (c *randCommand) GetMinRequiredRank() float64 {
	return c.AllowedRank
}

func (c *randCommand) Validate(args []string) error {
	return nil
}

func (c *randCommand) Exec(args []string) (*CommandResult, error) {
	max := 100
	min := 0
	random := rand.Intn(max-min) + min
	payloads := []*Event{
		{
			Method:  "chatMsg",
			Message: EventPayload{Message: strconv.Itoa(random), Meta: struct{}{}},
		},
	}
	return NewCommandResult(payloads), nil
}
