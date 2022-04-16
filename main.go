package main

import (
	"runtime"

	"go.uber.org/zap"

	"github.com/uninstallgentoo/go-syncbot/client"
	"github.com/uninstallgentoo/go-syncbot/command"
	"github.com/uninstallgentoo/go-syncbot/commands"
	"github.com/uninstallgentoo/go-syncbot/config"
	"github.com/uninstallgentoo/go-syncbot/processors"
	"github.com/uninstallgentoo/go-syncbot/repository"
	"github.com/uninstallgentoo/go-syncbot/storages"
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

	commandHandler := command.NewCommandHandler(botProcessors, logger)

	commandHandler.RegisterCommands(commands.Dice)

	c := client.NewSocketClient(conf, botProcessors.Chat, commandHandler, logger)
	c.Start()
}
