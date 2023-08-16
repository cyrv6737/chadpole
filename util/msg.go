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
}
