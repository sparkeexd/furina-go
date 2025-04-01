package main

import (
	"github.com/sparkeexd/mimo/commands/daily"
	"github.com/sparkeexd/mimo/commands/hello"
	"github.com/sparkeexd/mimo/internal/models"
)

func main() {
	commands := []map[string]models.Command{
		hello.Commands,
		daily.Commands,
	}

	jobs := []models.CronJob{}

	bot := NewBot()
	bot.Start(commands, jobs)
}
