package processors

import (
	"strings"
	"github.com/uninstallgentoo/go-syncbot/models"
)

type Executor interface {
	Validate(args []string) error
	Exec(args []string) (*CommandResult, error)
}

type CommandHandler interface {
	Handle(msg models.Message)
	Parse(text string) (*Command, bool)
	GetHandler(command string) Executor
	Execute(command string, args []string, userRank float64) *CommandResult
}

type commandHandler struct {
	processors *Processors
	commandResults chan Event
	commandList map[string]Executor
}

func NewCommandHandler(processors *Processors) CommandHandler {
	commandList := map[string]Executor{
	}
	err := processors.Command.InitRanks()
	if err != nil {
		//TODO: pass logger and write error
	}
	return &commandHandler{
		processors:  processors,
		commandList: commandList,
		commandResults: make(chan Event),
	}
}

func (c *commandHandler) Handle(msg models.Message) {
	cleanedMessage := msg.Clean()
	result, isCommand := c.Parse(cleanedMessage.Text)
	if isCommand {
		users := c.processors.Chat.GetUsers()
		result := c.Execute(result.Expr, result.Args, users[msg.Username].Rank)
		if result != nil {
			for _, response := range result.Results {
				if response != nil {
					c.commandResults <- *response
				}
			}
		}
	}
}

func (c *commandHandler) GetHandler(command string) Executor {
	return c.commandList[command]
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
	if c.processors.Command.GetCommandRank(command) > userRank {
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

