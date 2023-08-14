/*
Example of button functions editing it's own message
*/
package commands

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

var page_number = 0

func RibbitBtnEditHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Ribbit ribbit",
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Prev",
							Style:    discordgo.PrimaryButton,
							CustomID: "rp_prev",
						},
						discordgo.Button{
							Label:    "Next",
							Style:    discordgo.PrimaryButton,
							CustomID: "rp_next",
						},
					},
				},
			},
		},
	})

	log.Println("[INFO] Sending ribbit with some prev and next buttons")
}

func RPPrevBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	page_number--
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You pressed Prev. %d", page_number),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Prev",
							Style:    discordgo.PrimaryButton,
							CustomID: "rp_prev",
						},
						discordgo.Button{
							Label:    "Next",
							Style:    discordgo.PrimaryButton,
							CustomID: "rp_next",
						},
					},
				},
			},
		},
	})

	log.Println("[INFO] User pressed prev")
}

func RPNextBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	page_number++
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("You pressed Prev. %d", page_number),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Prev",
							Style:    discordgo.PrimaryButton,
							CustomID: "rp_prev",
						},
						discordgo.Button{
							Label:    "Next",
							Style:    discordgo.PrimaryButton,
							CustomID: "rp_next",
						},
					},
				},
			},
		},
	})

	log.Println("[INFO] User pressed next")
}
