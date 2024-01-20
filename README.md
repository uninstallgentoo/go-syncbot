# Go-SyncBot

Go-SyncBot is a lightweight and extensible bot for [CyTube](https://cytu.be/) project.

Designed to enhance the functionality and user experience on CyTube via custom chat commands. 

## Prerequisites
* Golang 1.21
* Docker 19+ (optional)

## Installation

### 1. Download

Clone or download the repository from [https://github.com/uninstallgentoo/go-syncbot](https://github.com/uninstallgentoo/go-syncbot).

### 2. Configuration
Go-SyncBot uses a configuration file to define CyTube server settings. The default configuration file is named config.example.yaml.

```bash
cp config.example.yaml config.yaml
```

Adjust the values in the configuration file to match your setup.
Example config.yaml:

``` yaml
server:
  host: "cytu.be"
  secure: true
channel:
  name: "test"
  password: ""
user:
  name: "admin"
  password: "admin"
database:
  path: "storages/bot.db"
```

## Run via Docker
```bash
docker compose up
```


## Build

Apply migrations to the database using the following command:

```bash
# install goose as a database migration tool
go install github.com/pressly/goose/v3/cmd/goose@latest

goose -dir ./migrations sqlite3 ./storages/bot.db up
```

Navigate to the repository directory and build the executable using the following command:

```bash
go mod download
go build -o syncbot main.go
```

Once the configuration is set, run the syncbot executable to run bot:

``` bash
./syncbot
```

## How to create a new command
1. Create a new Go module in the `commands` directory.
2. There's a simple interface to implement your own command. Pointer to the command struct in the `ExecFunc` signature allows you to use any attribute of it, such as Processors, Config and Cache.

`echo.go content:`
```go
var Echo = &command.Command{
	Name:        "echo", // command name that will be used in the chat
	Description: "display text that is passed as arguments",
	Rank:        1, // minimum required user rank to execute the command
	ExecFunc: func(args []string, cmd *command.Command) (models.CommandResult, error) {
		return models.NewCommandResult(
			models.NewChatMessage(strings.Join(args, " ")),
		), nil
	},
}
```
3. Register a new command in the `main.go`.
```go
commandHandler.RegisterCommands(
    commands.Dice,
    commands.Alert,
    commands.Stat,
    commands.Pick,
    commands.MagicBall,
    commands.Who,
    commands.Weather,
    commands.Echo,
)
```
4. Run the bot and send `!echo` in the CyTube channel chat.

## License
Go-SyncBot is open-source and distributed under the MIT License. Feel free to use, modify, and distribute it as per the terms of the license.

