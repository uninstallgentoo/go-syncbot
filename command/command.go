package command

import (
	"strings"

	"go.uber.org/zap"

	"github.com/uninstallgentoo/go-syncbot/config"
	"github.com/uninstallgentoo/go-syncbot/models"
	"github.com/uninstallgentoo/go-syncbot/processors"
	"github.com/uninstallgentoo/go-syncbot/storages"
)

type CommandExecutor interface {
	GetCommand() *Command
	GetValidateFunc() func([]string) error
	GetName() string
	GetDescription() string
	GetRank() float64
	SetProcessors(processors.Processors)
	SetConfig(*config.Config)
	SetCache(*storages.CacheStorage)
	Exec(args []string, cmd *Command) (models.CommandResult, error)
	Validate(args []string) error
}

type Command struct {
	Name         string
	Description  string
	Rank         float64
	ExecFunc     func([]string, *Command) (models.CommandResult, error)
	Processors   processors.Processors
	Config       *config.Config
	Cache        *storages.CacheStorage
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

func (c *Command) SetConfig(conf *config.Config) {
	c.Config = conf
}

func (c *Command) SetCache(cache *storages.CacheStorage) {
	c.Cache = cache
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
	cache          *storages.CacheStorage
	commandResults chan models.Event
	commandList    map[string]CommandExecutor
	logger         *zap.Logger
	conf           *config.Config
}

func NewCommandHandler(processors processors.Processors, cache *storages.CacheStorage, logger *zap.Logger, conf *config.Config) Handler {
	commandList := map[string]CommandExecutor{}
	err := processors.Command.InitRanks()
	if err != nil {
		logger.Error("init ranks failed", zap.Error(err))
	}
	return &commandHandler{
		processors:     processors,
		cache:          cache,
		commandList:    commandList,
		commandResults: make(chan models.Event),
		conf:           conf,
	}
}

func (c *commandHandler) GetCommandResults() chan models.Event {
	return c.commandResults
}

func (c *commandHandler) RegisterCommands(commands ...CommandExecutor) {
	for _, cmd := range commands {
		cmd.SetProcessors(c.processors)
		cmd.SetConfig(c.conf)
		cmd.SetCache(c.cache)
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
				if response.Message != nil {
					c.commandResults <- response
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
		return models.NewCommandResult(
			models.NewChatMessage("Permission denied for execution the command."),
		)
	}

	if command.GetValidateFunc() != nil {
		if err := command.Validate(args); err != nil {
			return models.NewCommandResult(
				models.NewChatMessage(err.Error()),
			)
		}
	}

	result, err := command.Exec(args, command.GetCommand())
	if err != nil {
		c.logger.Error("Error has occurred during command execution", zap.Error(err))
		return models.NewCommandResult(
			models.NewChatMessage(err.Error()),
		)

	}
	return result
}
