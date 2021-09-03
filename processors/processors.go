package processors

import (
	"sync-bot/repository"
)

type Processors struct {
	Command CommandProcessor
	Chat    ChatHandler
}

func NewProcessors(repo *repository.Repositories) *Processors {
	return &Processors{
		Command: NewCommandProcessor(repo.Command),
		Chat:    NewChatHandler(repo),
	}
}
