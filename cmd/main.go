package main

import (
	"github.com/sparkeexd/mimo/commands/daily"
	"github.com/sparkeexd/mimo/commands/hello"
	"github.com/sparkeexd/mimo/internal/models"
)

func main() {
	bot := models.NewBot(
		hello.Commands,
		daily.Commands,
	)

	bot.Start()
}
