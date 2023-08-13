/*
Convert link to odesli song.link
*/
package commands

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/bwmarrin/discordgo"
)

type OdesliResponse struct {
	Status      int    `json:"statusCode"`
	SongLinkURL string `json:"pageUrl"`
}

var (
	OdesliOptions = []*discordgo.ApplicationCommandOption{
		{
			Name:        "link",
			Description: "The music link to be converted",
			Type:        discordgo.ApplicationCommandOptionString,
		},
	}
)

func OdesliHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if len(i.ApplicationCommandData().Options) > 0 { // Make sure that the user actually gave a string
		link := i.ApplicationCommandData().Options[0].StringValue()
		log.Println("[INFO] Using link: " + link)
		message_content := getSongLink(link)
		log.Println("[INFO] Sending song link: " + message_content)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: message_content,
			},
		})
	} else {
		log.Println("[ERROR] No link provided by user for Odesli command")
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No link provided",
			},
		})
	}

}

func getSongLink(user_url string) string {
	api_url := "https://api.song.link/v1-alpha.1/links"
	user_url = url.QueryEscape(user_url)
	request_url := api_url + "?url=" + user_url

	response, err := http.Get(request_url)
	if err != nil {
		log.Println("Bad URL")
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("[ERROR] Couldn't read json body: ", err)
	}

	var result OdesliResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("[ERROR] Could not decode Odesli API response: ", err)
	}
	if result.Status == 400 {
		log.Println("[ERROR] Bad link sent to Odesli API")
		return "Bad Link"
	} else {
		odesli_link := result.SongLinkURL
		return odesli_link
	}

}
