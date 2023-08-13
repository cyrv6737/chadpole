package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"chadpole/commands"

	"github.com/bwmarrin/discordgo"
)

/*
Declare a list of commands and command handlers here. The functionality
of the handlers will be in the 'commands' package
*/
var (
	commandsList = []*discordgo.ApplicationCommand{
		{
			Name:        "ribbit",
			Description: "Ribbit",
		},
		{
			Name:        "ribbit-embed",
			Description: "Ribbit, but embeded",
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ribbit":       commands.RibbitHandler,
		"ribbit-embed": commands.RibbitEmbedHandler,
	}
)

func main() {

	bot_token := os.Getenv("DC_TOKEN") // Take bot token as an env var

	// Create the bot instance
	discord, err := discordgo.New("Bot " + bot_token)
	if err != nil {
		log.Fatal(err)
	}

	// Open up the bot instance, defer close for when bot is interrupted
	discord.Open()
	defer discord.Close()

	// Register commands
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commandsList))
	for i, v := range commandsList {
		cmd, err := discord.ApplicationCommandCreate(discord.State.User.ID, discord.State.Application.GuildID, v)
		if err != nil {
			log.Fatal(err)
		}
		registeredCommands[i] = cmd
	}

	// Attach handlers to commands
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})

	// Run until terminated in the console
	fmt.Println("Chadpole is ribbiting...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
