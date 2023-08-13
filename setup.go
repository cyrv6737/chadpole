package main

import (
	"chadpole/commands"
	"log"

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

// Register all commands for the bot
func RegisterAllCommands(s *discordgo.Session) {
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commandsList))
	for i, v := range commandsList {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, s.State.Application.GuildID, v)
		if err != nil {
			log.Fatal(err)
		}
		registeredCommands[i] = cmd
	}
}

// Attach all of the handlers required for functionality
func SetupAllHandlers(s *discordgo.Session) {

	// Attach type-specific handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
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
	s.AddHandler(MessageCreateHandler)
}
