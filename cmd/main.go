package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/sparkeexd/mimo/commands/daily"
	"github.com/sparkeexd/mimo/commands/hello"
)

func main() {
	log.Println("Creating discord bot session...")
	CreateSession()

	log.Println("Adding commands...")
	AddCommands(hello.Commands)
	AddCommands(daily.Commands)

	// Event listener to stop the bot.
	log.Println("Bot is now running! Press Ctrl+C to exit.")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("Closing discord bot session...")
	CloseSession()
}
