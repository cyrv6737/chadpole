package main

import (
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
	log.Println("[START] Chadpole is ribbiting...")
	defer chadpole.Close()

	RegisterAllCommands(chadpole)
	SetupAllHandlers(chadpole)

	// Run until terminated in the console
	log.Println("[INFO] Ready")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	log.Println("[SHUTDOWN] Unribbiting gracefully...")

}
