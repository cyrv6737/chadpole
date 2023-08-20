package util

import (
	"io"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/otiai10/gosseract/v2"
)

/*
Make the tesseract client, pull the attachment via URL and save it as bytes.
Run tesseract over the saved bytes and return OCR string
*/
func DiscordImageToOCR(a *discordgo.MessageAttachment) string {
	ocr_client := gosseract.NewClient()
	defer ocr_client.Close()

	image_url := a.URL

	response, err := http.Get(image_url)
	if err != nil {
		log.Println("[ERROR] Getting attachment")
	}
	defer response.Body.Close()

	image_bytes, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("[ERROR] Couldn't read image bytes")
	}

	ocr_client.SetImageFromBytes(image_bytes)
	text, _ := ocr_client.Text()

	log.Println("[INFO] Tesseract successfully scanned image")
	return text

}

/*
Builds OCR Response. Accepts either:
@chadpole ocr (with image attached)
OR
@chadpole ocr (in reply to another messagew ith an image attached)
*/
func OCRResponse(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Attachments) > 0 {
		for _, attachment := range m.Attachments {
			if CheckContentType(attachment.ContentType, imageTypes) {
				log.Println("[INFO] Sending OCR Message")
				s.ChannelMessageSendComplex(m.ChannelID, BuildOCRMessage(attachment, m, 0)) // Send the message reply
			} else {
				s.ChannelMessageSendComplex(m.ChannelID, BuildOCRMessage(attachment, m, 1))
			}
		}
	} else if m.MessageReference != nil {
		ref_message, err := s.ChannelMessage(m.MessageReference.ChannelID, m.MessageReference.MessageID)
		if err != nil {
			log.Println("[ERROR] Could not get referenced message")
		}
		if len(ref_message.Attachments) > 0 {
			for _, attachment := range ref_message.Attachments {
				if CheckContentType(attachment.ContentType, imageTypes) {
					log.Println("[INFO] Sending OCR Message")
					s.ChannelMessageSendComplex(m.ChannelID, BuildOCRMessage(attachment, m, 0)) // Send the message reply
				} else {
					s.ChannelMessageSendComplex(m.ChannelID, BuildOCRMessage(attachment, m, 1))
				}
			}
		} else {
			log.Println("[ERROR] No attachment found in OCR request")
			s.ChannelMessageSendComplex(m.ChannelID, BuildNoAttachMessage(m))
		}
	} else {
		log.Println("[ERROR] No attachment found in OCR request")
		s.ChannelMessageSendComplex(m.ChannelID, BuildNoAttachMessage(m))
	}
}

func BuildOCRMessage(a *discordgo.MessageAttachment, m *discordgo.MessageCreate, c int) *discordgo.MessageSend {
	var msg_content string

	switch c {
	case 0:
		msg_content = "```\n" + DiscordImageToOCR(a) + "\n```"
	case 1:
		msg_content = "Invalid attachment content type"
	}

	message := discordgo.MessageSend{
		Content: msg_content,
		Reference: &discordgo.MessageReference{ // Reply to the message in question
			MessageID: m.ID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		},
	}

	return &message
}

func BuildNoAttachMessage(m *discordgo.MessageCreate) *discordgo.MessageSend {
	message := discordgo.MessageSend{
		Content: "No attachments found in message",
		Reference: &discordgo.MessageReference{ // Reply to the message in question
			MessageID: m.ID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		},
	}

	return &message
}
