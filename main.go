package main

import (
	"runtime"

	"go.uber.org/zap"

	"sync-bot/pkg/client"
	"sync-bot/pkg/commands"
	"sync-bot/pkg/config"
	"sync-bot/pkg/repository"
	"sync-bot/pkg/service"
	"sync-bot/pkg/storages"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	logger := zap.NewExample().Sugar()
	defer logger.Sync()
	conf, err := config.NewConfig()
	if err != nil {
		logger.Fatalf("Error parsing config: %e", err)
	}
	db := storages.NewDBConnection(conf)
	repo := repository.NewSyncRepository(db)
	chatService := service.NewService(repo)

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)

	commandHandler := commands.NewCommandHandler(chatService)
	chatHandler := client.NewChatHandler(chatService, userService, commandHandler)
	c := client.NewSocketClient(conf, chatHandler, logger)
	c.Start()
}
