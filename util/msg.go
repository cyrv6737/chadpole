package util

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func MessageCreateHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID { //Skip if message is from this bot
		return
	}

	if strings.Contains(strings.ToLower(m.Content), "frog") || strings.Contains(strings.ToLower(m.Content), "ribbit") {
		s.MessageReactionAdd(m.ChannelID, m.ID, "🐸") // Add frog emoji reaction
		log.Println("[INFO] Frog reacting to message " + m.ID)
	}

	if strings.Contains(strings.ToLower(m.Content), "titanfall 3") {
		roleID := "1135817463108476929"
		permission := int64(discordgo.PermissionSendMessages)
		/*
			Arg 1: Channel ID, get this from the created message
			Arg 2: Target ID, either memberID or role ID. Get this useing dev mode in discord client
			Arg 3: Type of Perm override. 0 for role 1 for member.
			Arg 4: "Allow" field. If you want to allow a permssion, set this field to the permission int64 cast.
			Arg 5: "Deny" field. If you want to deny a permssion, set this field to the permission int64 cast.
			Whichever of the two between arg 4 and 5 is not used, fill in a 0
			Permission must be cast to an int64
			If I wanted to allow the permission, I would do the following:
			s.ChannelPermissionSet(m.ChannelID, roleID, 0, permission, 0)
		*/
		s.ChannelPermissionSet(m.ChannelID, roleID, 0, 0, permission)
	}

	/*
		Check for mentions to the bot itself in the message
	*/
	if m.Mentions != nil {
		for _, mention := range m.Mentions {
			if mention.ID == s.State.User.ID {
				if strings.Contains(strings.ToLower(m.Content), "ocr") { // @chadpole ocr
					log.Println("[INFO] Received OCR Request")
					go OCRResponse(s, m) // Start running OCR process on its own goroutine
				}
			}
		}
	}

}
