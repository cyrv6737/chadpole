package util

import (
	"io"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/otiai10/gosseract/v2"
)

var (
	imageTypes = []string{"image/png", "image/jpeg", "image/bmp", "image/webp", "image/tiff"}
)

/*
Check to make sure content type is valid for tesseract
*/
func CheckContentType(c string, l []string) bool {
	for _, item := range l {
		if item == c {
			log.Println("[INFO] Found valid ContentType for OCR")
			return true
		}
	}
	log.Println("[ERROR] Invalid ContentType for OCR")
	return false
}

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
				s.ChannelMessageSendComplex(m.ChannelID, BuildOCRMessage(attachment, m)) // Send the message reply
			} else {
				the_message := discordgo.MessageSend{
					Content: "Invalid attachment content type",
					Reference: &discordgo.MessageReference{ // Reply to the message in question
						MessageID: m.ID,
						ChannelID: m.ChannelID,
						GuildID:   m.GuildID,
					},
				}
				s.ChannelMessageSendComplex(m.ChannelID, &the_message) // Send the message reply
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
					s.ChannelMessageSendComplex(m.ChannelID, BuildOCRMessage(attachment, m)) // Send the message reply
				} else {
					the_message := discordgo.MessageSend{
						Content: "Invalid attachment content type",
						Reference: &discordgo.MessageReference{ // Reply to the message in question
							MessageID: m.ID,
							ChannelID: m.ChannelID,
							GuildID:   m.GuildID,
						},
					}
					s.ChannelMessageSendComplex(m.ChannelID, &the_message) // Send the message reply
				}
			}
		} else {
			log.Println("[ERROR] No attachment found in OCR request")
			the_message := discordgo.MessageSend{
				Content: "No attachment found",
				Reference: &discordgo.MessageReference{ // Reply to the message in question
					MessageID: m.ID,
					ChannelID: m.ChannelID,
					GuildID:   m.GuildID,
				},
			}
			s.ChannelMessageSendComplex(m.ChannelID, &the_message) // Send the message reply
		}
	} else {
		log.Println("[ERROR] No attachment found in OCR request")
		the_message := discordgo.MessageSend{
			Content: "No attachment found",
			Reference: &discordgo.MessageReference{ // Reply to the message in question
				MessageID: m.ID,
				ChannelID: m.ChannelID,
				GuildID:   m.GuildID,
			},
		}
		s.ChannelMessageSendComplex(m.ChannelID, &the_message) // Send the message reply
	}
}

func BuildOCRMessage(a *discordgo.MessageAttachment, m *discordgo.MessageCreate) *discordgo.MessageSend {
	the_message := discordgo.MessageSend{
		Content: "```\n" + DiscordImageToOCR(a) + "\n```",
		Reference: &discordgo.MessageReference{ // Reply to the message in question
			MessageID: m.ID,
			ChannelID: m.ChannelID,
			GuildID:   m.GuildID,
		},
	}

	return &the_message
}
