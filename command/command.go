package command

import (
	"strings"
	"sync"

	"github.com/uninstallgentoo/go-syncbot/models"
	"github.com/uninstallgentoo/go-syncbot/processors"
)

type Command struct {
	Name         string
	Description  string
	Rank         float64
	ExecFunc     func() (*models.CommandResult, error)
	ValidateFunc func([]string) error
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

func (c *Command) GetFunc() func() (*models.CommandResult, error) {
	return c.ExecFunc
}

func (c *Command) Validate() func([]string) error {
	return nil
}

type Executor interface {
	Validate(args []string) error
	Exec(args []string) (*models.CommandResult, error)
}

type Handler interface {
	Handle(msg models.Message)
	Parse(text string) (Command, []string, bool)
	Execute(command Command, args []string, userRank float64) *models.CommandResult
	RegisterCommands(commands ...Command)
	GetCommandResults() chan models.Event
}

type commandHandler struct {
	m              *sync.RWMutex
	processors     *processors.Processors
	commandResults chan models.Event
	commandList    map[string]Command
}

func NewCommandHandler(processors *processors.Processors) Handler {
	commandList := map[string]Command{}
	err := processors.Command.InitRanks()
	if err != nil {
		//TODO: pass logger and write error
	}
	return &commandHandler{
		m:              &sync.RWMutex{},
		processors:     processors,
		commandList:    commandList,
		commandResults: make(chan models.Event),
	}
}

func (c *commandHandler) GetCommandResults() chan models.Event {
	return c.commandResults
}

func (c *commandHandler) RegisterCommands(commands ...Command) {
	c.m.RLock()
	defer c.m.RUnlock()
	for _, cmd := range commands {
		c.commandList[cmd.GetName()] = cmd
	}
}

func (c *commandHandler) Handle(msg models.Message) {
	cleanedMessage := msg.Clean()
	command, args, isCommand := c.Parse(cleanedMessage.Text)
	if isCommand {
		users := c.processors.Chat.GetUsers()
		result := c.Execute(command, args, users[msg.Username].Rank)
		if result != nil {
			for _, response := range result.Results {
				if response != nil {
					c.commandResults <- *response
				}
			}
		}
	}
}

func (c *commandHandler) GetHandler(command string) Command {
	return c.commandList[command]
}

func (c *commandHandler) Parse(text string) (Command, []string, bool) {
	if strings.HasPrefix(text, "!") {
		message := strings.Split(text, " ")
		cmd, args := message[0], message[1:]
		if command, ok := c.commandList[cmd]; ok {
			return command, args, true
		}
	}
	return Command{}, nil, false
}

func (c *commandHandler) Execute(command Command, args []string, userRank float64) *models.CommandResult {
	if command.GetRank() > userRank {
		return processors.NewCommandResult([]*models.Event{
			{
				Method:  "chatMsg",
				Message: models.EventPayload{Message: "Permission denied for execution the command.", Meta: struct{}{}},
			},
		})
	}
	err := command.Validate()(args)
	if err != nil {
		return processors.NewCommandResult([]*models.Event{
			{
				Method:  "chatMsg",
				Message: models.EventPayload{Message: err.Error(), Meta: struct{}{}},
			},
		})
	}
	result, err := command.GetFunc()()
	if err != nil {
		return processors.NewCommandResult([]*models.Event{
			{
				Method:  "chatMsg",
				Message: models.EventPayload{Message: err.Error(), Meta: struct{}{}},
			},
		})
	}
	return result
}
