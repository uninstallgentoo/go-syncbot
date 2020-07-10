package commands

import (
	"strings"

	"sync-bot/pkg/service"
)

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

type CommandExecutor interface {
	GetMinRequiredRank() float64
	Validate(args []string) error
	Exec(args []string) (*CommandResult, error)
}

type CommandHandler interface {
	Parse(text string) (*Command, bool)
	GetHandler(command string) CommandExecutor
	Execute(command string, args []string, userRank float64) *CommandResult
}

type Command struct {
	Expr string
	Args []string
}

type commandHandler struct {
	chatService service.Chat
	commandList map[string]CommandExecutor
}

func NewCommandHandler(s service.Chat) CommandHandler {
	commandList := map[string]CommandExecutor{
		"rand":   NewRandCommand(),
		"add":    NewAddCommand(),
		"random": NewRandomMessageCommand(s),
		"4chan":  NewFourchanCommand(),
	}
	return &commandHandler{
		chatService: s,
		commandList: commandList,
	}
}

func (c *commandHandler) GetHandler(command string) CommandExecutor {
	return c.commandList[command]
}

func (c *commandHandler) IsCommandExecAllowed(userRank float64, handler CommandExecutor) bool {
	if userRank < handler.GetMinRequiredRank() {
		return false
	}
	return true
}

func (c *commandHandler) Parse(text string) (*Command, bool) {
	splittedMessageText := strings.Split(text, " ")
	expr := splittedMessageText[0]
	if strings.HasPrefix(expr, "!") {
		args := splittedMessageText[1:]
		command := &Command{
			//remove prefix from command
			Expr: expr[1:],
			Args: args,
		}
		return command, true
	}
	return nil, false
}

func (c *commandHandler) Execute(command string, args []string, userRank float64) *CommandResult {
	handler := c.GetHandler(command)
	if handler == nil {
		return NewCommandResult([]*Event{
			{
				Method:  "chatMsg",
				Message: EventPayload{Message: "Команда отсутствует.", Meta: struct{}{}},
			},
		})
	}
	if !c.IsCommandExecAllowed(userRank, handler) {
		return NewCommandResult([]*Event{
			{
				Method:  "chatMsg",
				Message: EventPayload{Message: "Отсутствуют права для выполнения команды.", Meta: struct{}{}},
			},
		})
	}
	err := handler.Validate(args)
	if err != nil {
		return NewCommandResult([]*Event{
			{
				Method:  "chatMsg",
				Message: EventPayload{Message: err.Error(), Meta: struct{}{}},
			},
		})
	}
	result, err := handler.Exec(args)
	if err != nil {
		return NewCommandResult([]*Event{
			{
				Method:  "chatMsg",
				Message: EventPayload{Message: err.Error(), Meta: struct{}{}},
			},
		})
	}
	return result
}
