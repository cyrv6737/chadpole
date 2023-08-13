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
		{
			Name:        "ribbit-button",
			Description: "Ribbit with buttons",
		},
		{
			Name:        "odesli",
			Description: "Convert any music link into a song.link",
			Options:     commands.OdesliOptions,
		},
	}
	commandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"ribbit":        commands.RibbitHandler,
		"ribbit-embed":  commands.RibbitEmbedHandler,
		"ribbit-button": commands.RibbitButtonHandler,
		"odesli":        commands.OdesliHandler,
	}
	componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"primary_test": commands.PrimaryTestBtnHandler,
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
	/*
		// Attach command handlers
		discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		})
	*/

	// Attach all handlers
	discord.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// Attach command handlers for slash commands
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		// Attach component handlers, such as handlers for buttons
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	// Allow the created message handler to monitor messages
	discord.AddHandler(MessageCreateHandler)

	// Run until terminated in the console
	fmt.Println("Chadpole is ribbiting...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

}
