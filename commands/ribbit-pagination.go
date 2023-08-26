/*
Exmaple pagination implementation with buttons.
Inspired from https://github.com/Necroforger/dgwidgets/blob/master/paginator.go
As well as from my own personal work at: https://github.com/CooldudePUGS/Spectre/blob/90463d95839caf6a8551cf6fa91ac2dc952101b5/cogs/ModSearch.py

Fetches a JSON response from the locally hosted frog information API, displays that information in a paginated embed.
*/
package commands

import (
	"chadpole/widgets"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

type FrogExample struct {
	FrogName  string `json:"name"`
	FrogDesc  string `json:"desc"`
	FrogLink  string `json:"link"`
	FrogImage string `json:"imageURL"`
}

/*
Entrypoint for the pagination system
*/
func RibbitPaginationHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Put API fetch on its own goroutine, incase server takes a while to respond
	// Not an issue with the local API but could be elsewhere
	jsonChannel := make(chan []byte)
	go FetchAPIData(jsonChannel)
	jsonResult := <-jsonChannel

	new_pagination := widgets.PaginationView{
		Data:                []widgets.PageData{}, // Declare slice of PageData struct to store page data
		EnableLink:          true,
		EnableStop:          false,
		EnableShowInChannel: true,
		EnableSecondRow:     true,
		IsEphemeral:         true,
	}

	var jsonData []FrogExample

	err := json.Unmarshal(jsonResult, &jsonData)
	if err != nil {
		log.Println("[ERROR] Could not decode json")
	}

	// Iterate over the json data, feed it into the paginator's PageData slice
	for _, jsonItem := range jsonData {
		onePage := widgets.PageData{
			Title:    jsonItem.FrogName,
			Desc:     jsonItem.FrogDesc,
			Link:     jsonItem.FrogLink,
			ImageURL: jsonItem.FrogImage,
		}
		new_pagination.Data = append(new_pagination.Data, onePage)
	}

	log.Println("[INFO] New pagination created")
	new_pagination.SendMessage(s, i) // Send the message, functions as the entrypoint for the pagination view
}

func FetchAPIData(ch chan []byte) {
	/*
		Fetches frog data from the API hosted locally by the bot
	*/
	response, err := http.Get("http://127.0.0.1:8081/frog")
	if err != nil {
		log.Println("[ERROR] Could not get API response")
		ch <- nil
	}
	defer response.Body.Close()

	jsonBody, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("[ERROR] could not read JSON body")
	}

	log.Println("[INFO] Successfully retrieved API response")
	ch <- jsonBody
}
