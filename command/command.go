package command

import (
	"go.uber.org/zap"
	"strings"

	"github.com/uninstallgentoo/go-syncbot/models"
	"github.com/uninstallgentoo/go-syncbot/processors"
)

type CommandExecutor interface {
	GetCommand() *Command
	GetValidateFunc() func([]string) error
	GetName() string
	GetDescription() string
	GetRank() float64
	SetProcessors(processors.Processors)
	Exec(args []string, cmd *Command) (models.CommandResult, error)
	Validate(args []string) error
}

type Command struct {
	Name         string
	Description  string
	Rank         float64
	ExecFunc     func([]string, *Command) (models.CommandResult, error)
	Processors   processors.Processors
	ValidateFunc func([]string) error
}

func (c *Command) GetCommand() *Command {
	return c
}

func (c *Command) GetValidateFunc() func([]string) error {
	return c.ValidateFunc
}

func (c *Command) GetName() string {
	return c.Name
}

func (c *Command) GetDescription() string {
	return c.Description
}

func (c *Command) GetRank() float64 {
	return c.Rank
}

func (c *Command) SetProcessors(p processors.Processors) {
	c.Processors = p
}

func (c *Command) Exec(args []string, cmd *Command) (models.CommandResult, error) {
	return c.ExecFunc(args, cmd)
}

func (c *Command) Validate(args []string) error {
	return c.ValidateFunc(args)
}

type Handler interface {
	Handle(msg models.Message)
	Parse(text string) (CommandExecutor, []string, bool)
	Execute(command CommandExecutor, args []string, userRank float64) models.CommandResult
	RegisterCommands(commands ...CommandExecutor)
	GetCommandResults() chan models.Event
}

type commandHandler struct {
	processors     processors.Processors
	commandResults chan models.Event
	commandList    map[string]CommandExecutor
	logger         *zap.Logger
}

func NewCommandHandler(processors processors.Processors, logger *zap.Logger) Handler {
	commandList := map[string]CommandExecutor{}
	err := processors.Command.InitRanks()
	if err != nil {
		logger.Error("init ranks failed", zap.Error(err))
	}
	return &commandHandler{
		processors:     processors,
		commandList:    commandList,
		commandResults: make(chan models.Event),
	}
}

func (c *commandHandler) GetCommandResults() chan models.Event {
	return c.commandResults
}

func (c *commandHandler) RegisterCommands(commands ...CommandExecutor) {
	for _, cmd := range commands {
		cmd.SetProcessors(c.processors)
		c.commandList[cmd.GetName()] = cmd
	}
}

func (c *commandHandler) Handle(msg models.Message) {
	cleanedMessage := msg.Clean()
	command, args, isCommand := c.Parse(cleanedMessage.Text)
	if isCommand {
		users := c.processors.Chat.GetUsers()
		result := c.Execute(command, args, users[msg.Username].Rank)
		if len(result.Results) > 0 {
			for _, response := range result.Results {
				if response != nil && response.Message != nil {
					c.commandResults <- *response
				}
			}
		}
	}
}

func (c *commandHandler) GetHandler(command string) CommandExecutor {
	return c.commandList[command]
}

func (c *commandHandler) Parse(text string) (CommandExecutor, []string, bool) {
	if strings.HasPrefix(text, "!") {
		message := strings.Split(text, " ")
		cmd, args := message[0][1:], message[1:]
		if command, ok := c.commandList[cmd]; ok {
			return command, args, true
		}
	}
	return nil, nil, false
}

func (c *commandHandler) Execute(command CommandExecutor, args []string, userRank float64) models.CommandResult {
	if command.GetRank() > userRank {
		return processors.NewCommandResult([]*models.Event{
			{
				Method:  "chatMsg",
				Message: models.EventPayload{Message: "Permission denied for execution the command.", Meta: struct{}{}},
			},
		})
	}

	if command.GetValidateFunc() != nil {
		if err := command.Validate(args); err != nil {
			return processors.NewCommandResult([]*models.Event{
				{
					Method:  "chatMsg",
					Message: models.EventPayload{Message: err.Error(), Meta: struct{}{}},
				},
			})
		}
	}

	result, err := command.Exec(args, command.GetCommand())
	if err != nil {
		c.logger.Error("Error has occurred during command execution", zap.Error(err))
		return processors.NewCommandResult([]*models.Event{
			{
				Method:  "chatMsg",
				Message: models.EventPayload{Message: err.Error(), Meta: struct{}{}},
			},
		})
	}
	return result
}
