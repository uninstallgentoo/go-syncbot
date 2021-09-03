package processors

import (
	"sync-bot/models"
	"sync-bot/repository"
)

type CommandProcessor interface {
	Add(command *models.Command) error
	UpdateRank(command *models.Command) error
	GetCommandRank(command string) float64
	InitRanks() error
}

type EventPayload struct {
	Message string      `json:"msg"`
	Meta    interface{} `json:"meta"`
}

type Event struct {
	Method  string
	Message interface{}
}

type CommandResult struct {
	Results []*Event
}

func NewCommandResult(results []*Event) *CommandResult {
	return &CommandResult{
		Results: results,
	}
}

type Command struct {
	Expr string
	Args []string
}

type commandProcessor struct {
	repo         repository.CommandRepository
	commandRanks map[string]float64
}

func NewCommandProcessor(repo repository.CommandRepository) CommandProcessor {
	return &commandProcessor{
		repo: repo,
		commandRanks: map[string]float64{
			"rand":       1,
			"add":        1,
			"random":     1,
			"4chan":      1,
			"setCmdRank": 3,
		},
	}
}

func (s *commandProcessor) Add(command *models.Command) error {
	return s.repo.Add(command)
}

func (s *commandProcessor) UpdateRank(command *models.Command) error {
	err := s.repo.UpdateRank(command)
	if err != nil {
		return err
	}
	s.commandRanks[command.Command] = command.Rank
	return nil
}

func (s *commandProcessor) GetCommandRank(command string) float64 {
	return s.commandRanks[command]
}

func (s *commandProcessor) InitRanks() error {
	commands, err := s.repo.FetchAll()
	if err != nil {
		return err
	}
	for _, item := range commands {
		s.commandRanks[item.Command] = item.Rank
	}
	return nil
}
