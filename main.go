package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
)

func main() {

	bot_token := os.Getenv("DC_TOKEN") // Take bot token as an env var

	// Create the bot instance
	chadpole, err := discordgo.New("Bot " + bot_token)
	if err != nil {
		log.Fatal(err)
	}

	// Open up the bot instance, defer close for when bot is interrupted
	chadpole.Open()
	defer chadpole.Close()

	RegisterAllCommands(chadpole)
	SetupAllHandlers(chadpole)

	// Run until terminated in the console
	fmt.Println("Chadpole is ribbiting...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
