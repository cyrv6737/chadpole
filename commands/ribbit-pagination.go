/*
Exmaple pagination implementation with buttons.
Inspired from https://github.com/Necroforger/dgwidgets/blob/master/paginator.go
As well as from my own personal work at: https://github.com/CooldudePUGS/Spectre/blob/90463d95839caf6a8551cf6fa91ac2dc952101b5/cogs/ModSearch.py

One notable difference between this implementation and the python one in spectre is that
I cannot properly get multiple indepdendent paginations up at once.
The values are properly isolated due to making a new struct each time as well as using mutex locks
I believe the issue lies in that the button handlers get added to the global bot
To my knowledge there isn't a way to add these handlers to a specific message instance. This might
be due to the fact that there is no "view" implementation like disordpy has.

A lot of variables here start with "TS" which stands for thunderstore, as this was originally
supposed to be an implementation of thunderstore mod searching. Might save that for another module
and get extra confusing lol.
*/
package commands

import (
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

/*
Create a global variable to keep track of if the pagination system is running or not.
Safeguards against the issues stated above
*/
var isRunning bool

/*
Need to create a class/struct for the pagination view.
Add in any variables you will need in here. Some notable things
would be variables to hold JSON values if you're displaying a search
from an API for example
*/
type PaginationView struct {
	sync.Mutex
	index      int
	embedtitle string
	embeddesc  string
}

/*
Create a function for the struct that checks the global var if the running status
*/
func (p *PaginationView) Running() bool {
	p.Lock()
	running := isRunning
	p.Unlock()
	return running
}

/*
Creates the embed. This function is called every time there is an update to the message
*/
func (p *PaginationView) CreateEmbed() []*discordgo.MessageEmbed {
	p.embedtitle = fmt.Sprintf("%d", p.index)
	p.embeddesc = fmt.Sprintf("%d", p.index)

	embed := []*discordgo.MessageEmbed{
		{
			Title:       p.embedtitle,
			Description: p.embeddesc,
		},
	}

	return embed
}

/*
Function is called once at the start to create the first instance of the pagination.
Importantly, there is a check to make sure that there are no running paginations currently

Handers are added several times because currently they are set to add once to avoid duplication
*/
func (p *PaginationView) SendMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if p.Running() {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Pagination already running",
			},
		})
		log.Println("[ERROR] Pagination already running")
		p.TSAddHandlers(s, i)
		return
	}
	p.TSAddHandlers(s, i)
	// Using mutex locks just to be safe, even though realistically I don't have to
	p.Lock()
	isRunning = true
	p.Unlock()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: p.CreateEmbed(),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Prev",
							Style:    discordgo.SuccessButton,
							CustomID: "ts_prev",
						},
						discordgo.Button{
							Label:    "Next",
							Style:    discordgo.SuccessButton,
							CustomID: "ts_next",
						},
						discordgo.Button{
							Label:    "Done",
							Style:    discordgo.PrimaryButton,
							CustomID: "ts_done",
						},
						discordgo.Button{
							Label:    "Stop",
							Style:    discordgo.DangerButton,
							CustomID: "ts_stop",
						},
					},
				},
			},
		},
	})
	log.Println("[INFO] Sent initial pagination view")
}

/*
Updates the message every time next or prev is pressed
*/
func (p *PaginationView) UpdateMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.TSAddHandlers(s, i)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: p.CreateEmbed(),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label:    "Prev",
							Style:    discordgo.SuccessButton,
							CustomID: "ts_prev",
						},
						discordgo.Button{
							Label:    "Next",
							Style:    discordgo.SuccessButton,
							CustomID: "ts_next",
						},
						discordgo.Button{
							Label:    "Done",
							Style:    discordgo.PrimaryButton,
							CustomID: "ts_done",
						},
						discordgo.Button{
							Label:    "Stop",
							Style:    discordgo.DangerButton,
							CustomID: "ts_stop",
						},
					},
				},
			},
		},
	})
	log.Println("[INFO] Updated pagination")
}

/*
Increments the index value. In a proper situation you would use this index value to get more data from
a data structure
*/
func (p *PaginationView) TSNextBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	defer p.Unlock()
	p.index++
	log.Println("[INFO] Pagination data incremented")
	p.UpdateMessage(s, i)
}

func (p *PaginationView) TSPrevBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	defer p.Unlock()
	p.index--
	log.Println("[INFO] Pagination data decremented")
	p.UpdateMessage(s, i)
}

/*
Deletes the message if pressed, also resets running status
*/
func (p *PaginationView) TSStopBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	isRunning = false
	p.Unlock()
	/*
		Holy fucking shit
		this is how you delete your own fucking message
		why did it take so long to figure this out?
		fucks sake

		Do note: you can get the Channel ID of the Interaction with i.Message.ChannelID
		and it's own message id with i.Message.ID

	*/
	log.Println("[INFO] Pagination stopped")
	s.ChannelMessageDelete(i.Message.ChannelID, i.Message.ID)
	log.Println("[INFO] Pagination message removed from channel")
}

/*
Resets running status but instead of deleting the message, it just gets rid of the buttons so the
embed is permanent in the channel
*/
func (p *PaginationView) TSDoneBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	isRunning = false
	p.Unlock()
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseUpdateMessage,
		Data: &discordgo.InteractionResponseData{
			Embeds: p.CreateEmbed(),
		},
	})
	log.Println("[INFO] Pagination stopped, embed remains in channel")
}

/*
Add all the handlers for the buttons. Notably this is a direct copy-paste implementation from setup.go
Definitely not the best solution. Should probably change this later.
I am fairly confident (lol) that cleaner code here would not help the issue with concurrent paginations.
*/
func (p *PaginationView) TSAddHandlers(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var (
		componentHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
			"ts_next": p.TSNextBtnHandler,
			"ts_prev": p.TSPrevBtnHandler,
			"ts_stop": p.TSStopBtnHandler,
			"ts_done": p.TSDoneBtnHandler,
		}
	)
	s.AddHandlerOnce(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// Attach component handlers, such as handlers for buttons
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	log.Println("[INFO] Added pagination handlers")
}

/*
Entrypoint for the pagination system
*/
func RibbitPaginationHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	new_pagination := PaginationView{
		index: 0,
	}
	log.Println("[INFO] New pagination created")
	new_pagination.SendMessage(s, i)

}
