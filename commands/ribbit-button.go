/*
Basic response slash command example, but with some buttons?
*/
package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func RibbitButtonHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ribbit ribbit",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Test",
							Style: discordgo.LinkButton,
							URL:   "https://www.google.com",
						},
						discordgo.Button{
							Label:    "Primary",
							Style:    discordgo.PrimaryButton,
							CustomID: "primary_test",
						},
					},
				},
			},
		},
	})

	log.Println("[INFO] Sending ribbit with buttons")
}

func PrimaryTestBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ribbit, you pressed the primary button",
		},
	})

	log.Println("[INFO] User pressed the Primary Test button")
}
