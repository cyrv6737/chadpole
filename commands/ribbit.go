/*
Basic response slash command example
*/
package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func RibbitHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ribbit ribbit",
		},
	})

	log.Println("[INFO] Sending normal ribbit")
}
