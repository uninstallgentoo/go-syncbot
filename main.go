package main

import (
	"runtime"

	"go.uber.org/zap"

	"sync-bot/client"
	"sync-bot/config"
	"sync-bot/processors"
	"sync-bot/repository"
	"sync-bot/storages"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	logger, err := zap.NewProduction()
	defer logger.Sync()

	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatal("Error parsing config", zap.Error(err))
	}
	db := storages.NewDBConnection(conf)

	repositories := repository.NewRepositories(db)
	botProcessors := processors.NewProcessors(repositories)

	commandHandler := processors.NewCommandHandler(botProcessors)

	c := client.NewSocketClient(conf, botProcessors.Chat, commandHandler, logger)
	c.Start()
}
