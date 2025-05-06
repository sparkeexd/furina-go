package main

import (
	"github.com/sparkeexd/furina/internal/application/bot"
)

func main() {
	bot := bot.NewBot()
	bot.Start()
}
