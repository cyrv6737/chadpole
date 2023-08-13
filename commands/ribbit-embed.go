/*
Basic response slash command example, but ultizes an embedded message.
*/
package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func RibbitEmbedHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Ribbit Test",
					Description: "Embedded ribbiting",
					Thumbnail: &discordgo.MessageEmbedThumbnail{
						URL: "https://em-content.zobj.net/thumbs/160/twitter/53/frog-face_1f438.png",
					},
				},
			},
		},
	})

	log.Println("[INFO] Sending embedded ribbit")
}
