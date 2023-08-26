/*
Exmaple pagination implementation with buttons.
Inspired from https://github.com/Necroforger/dgwidgets/blob/master/paginator.go
As well as from my own personal work at: https://github.com/CooldudePUGS/Spectre/blob/90463d95839caf6a8551cf6fa91ac2dc952101b5/cogs/ModSearch.py
KNOWN ISSUES:
  - Logging information will repeat several times depending on how many times pagination has been called.
    Actual functionality of this is not affected since we make sure there are not multiple paginations running at once
    However this is pretty far from ideal
  - The above has been solved by generating a random 8 character prefix for the handler CustomIDs.
    This will ensure that different handlers are assigned to the pagination every time.
    The drawback of this is that there are now essentially "dead handlers" attached to the bot. Maybe
    the garbage collector deals with it at some point. The fuck do I know. At least it works now
*/

package widgets

import (
	"crypto/rand"
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

/*
Need to create a class/struct for the pagination view.
Add in any variables you will need in here. Some notable things
would be variables to hold JSON values if you're displaying a search
from an API for example
*/
type PaginationView struct {
	sync.Mutex
	Index               int
	EmbedTitle          string
	EmbedDesc           string
	EmbedLink           string
	EmbedImgURL         string
	handerPrefix        string
	PageBtnHandlers     map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	Data                []PageData // Create a slice of the struct where our data will be stored.
	currentPage         int
	EnableLink          bool
	IsEphemeral         bool
	EnableStop          bool
	EnableShowInChannel bool
	EnableSecondRow     bool
}

type PageData struct {
	Title    string
	Desc     string
	Link     string
	ImageURL string
}

/*
General setup housekeeping should go here
*/
func (p *PaginationView) Setup(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	p.handerPrefix = p.GenPrefix() // Generate prefix to uniquely identify paginator controls
	p.Unlock()
	// Assign handlers to their respective CustomIDs
	p.PageBtnHandlers[p.handerPrefix+"pg_next"] = p.PG_NextBtnHandler
	p.PageBtnHandlers[p.handerPrefix+"pg_prev"] = p.PG_PrevBtnHandler
	p.PageBtnHandlers[p.handerPrefix+"pg_stop"] = p.PG_StopBtnHandler
	p.PageBtnHandlers[p.handerPrefix+"pg_done"] = p.PG_DoneBtnHandler
	p.PageBtnHandlers[p.handerPrefix+"pg_first"] = p.PG_FirstBtnHandler
	p.PageBtnHandlers[p.handerPrefix+"pg_last"] = p.PG_LastBtnHandler
	// Add the handlers to the bot
	p.PG_AddHandlers(s, i)
	p.currentPage = p.Index + 1
}

/*
Generate a unique random prefix to identify a unique paginator's button CustomIDs
Unfortunately this is how I decided to work around being unable to cleanly add
buttons to a specific instance rather than globally
*/
func (p *PaginationView) GenPrefix() string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 8
	prefix := make([]byte, length)
	_, err := rand.Read(prefix)
	if err != nil {
		log.Println("[WARN] Could not generate random prefix. Things might get weird.")
		return ""
	}
	for i, b := range prefix {
		prefix[i] = charset[b%byte(len(charset))]
	}
	log.Printf("[INFO] Using prefix %s", prefix)
	return string(prefix)
}

/*
Creates the embed. This function is called every time there is an update to the message
*/
func (p *PaginationView) CreateEmbed() []*discordgo.MessageEmbed {
	// Pull data from Data slice based on the index which is modified with the buttons
	p.EmbedTitle = p.Data[p.Index].Title
	p.EmbedDesc = p.Data[p.Index].Desc
	p.EmbedLink = p.Data[p.Index].Link
	p.EmbedImgURL = p.Data[p.Index].ImageURL

	embed := []*discordgo.MessageEmbed{
		{
			Title:       p.EmbedTitle,
			Description: p.EmbedDesc,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf("Page: %d/%d", p.currentPage, len(p.Data)),
			},
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: p.EmbedImgURL,
			},
		},
	}

	return embed
}

/*
Since we need to create these buttons multiple times in the code, throw them in a function to improve readability
*/
func (p *PaginationView) CreateBtns() []discordgo.MessageComponent {

	var buttonComplex []discordgo.MessageComponent
	var secondRowBtns []discordgo.MessageComponent

	buttonRowOne := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "<<",
				Style:    discordgo.PrimaryButton,
				CustomID: p.handerPrefix + "pg_first",
			},
			discordgo.Button{
				Label:    "<",
				Style:    discordgo.PrimaryButton,
				CustomID: p.handerPrefix + "pg_prev",
			},
			discordgo.Button{
				Label:    ">",
				Style:    discordgo.PrimaryButton,
				CustomID: p.handerPrefix + "pg_next",
			},
			discordgo.Button{
				Label:    ">>",
				Style:    discordgo.PrimaryButton,
				CustomID: p.handerPrefix + "pg_last",
			},
		},
	}

	buttonComplex = append(buttonComplex, buttonRowOne)

	if p.EnableSecondRow {
		if p.EnableShowInChannel {
			secondRowBtns = append(secondRowBtns, discordgo.Button{
				Label:    "Show in Channel",
				Style:    discordgo.SuccessButton,
				CustomID: p.handerPrefix + "pg_done",
			})
		}

		if p.EnableLink {
			secondRowBtns = append(secondRowBtns, discordgo.Button{
				Label: "View",
				Style: discordgo.LinkButton,
				URL:   p.EmbedLink,
			})
		}

		if p.EnableStop {
			secondRowBtns = append(secondRowBtns, discordgo.Button{
				Label:    "Stop",
				Style:    discordgo.DangerButton,
				CustomID: p.handerPrefix + "pg_stop",
			})
		}
	} else {

		return buttonComplex

	}

	buttonRowTwo := discordgo.ActionsRow{
		Components: secondRowBtns,
	}

	buttonComplex = append(buttonComplex, buttonRowTwo)

	return buttonComplex

}

/*
Function is called once at the start to create the first instance of the pagination.
Importantly, there is a check to make sure that there are no running paginations currently

Handers are added several times because currently they are set to add once to avoid duplication
*/
func (p *PaginationView) SendMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Setup(s, i)               // Call setup function first so the handler prefix can be generated and all handlers added
	p.currentPage = p.Index + 1 // Display page number normally

	if p.IsEphemeral {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				// You can make searches only visible to the invoker with the following:
				// (Note: The stop button won't work)
				Flags:      discordgo.MessageFlagsEphemeral,
				Embeds:     p.CreateEmbed(),
				Components: p.CreateBtns(),
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds:     p.CreateEmbed(),
				Components: p.CreateBtns(),
			},
		})
	}

	log.Printf("[INFO] Sent initial pagination view for %s", p.handerPrefix)
}

/*
Updates the message every time next or prev is pressed
*/
func (p *PaginationView) UpdateMessage(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.currentPage = p.Index + 1 // Display page number normally
	if p.IsEphemeral {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				// You can make searches only visible to the invoker with the following:
				// (Note: The stop button won't work)
				Flags:      discordgo.MessageFlagsEphemeral,
				Embeds:     p.CreateEmbed(),
				Components: p.CreateBtns(),
			},
		})
	} else {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseUpdateMessage,
			Data: &discordgo.InteractionResponseData{
				Embeds:     p.CreateEmbed(),
				Components: p.CreateBtns(),
			},
		})
	}

	log.Printf("[INFO] Updated pagination %s", p.handerPrefix)
}

/*
Increments the index value. In a proper situation you would use this index value to get more data from
a data structure
*/
func (p *PaginationView) PG_NextBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	defer p.Unlock()
	p.Index = (p.Index + 1) % len(p.Data) // Circular pagination
	log.Printf("[INFO] Pagination %s data incremented", p.handerPrefix)
	p.UpdateMessage(s, i)
}

func (p *PaginationView) PG_PrevBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p.Lock()
	defer p.Unlock()
	if p.Index == 0 { // Prevent running out of bounds. Function as a "last" button if index is at 0
		p.Index = len(p.Data) - 1
		p.UpdateMessage(s, i)
	} else {
		p.Index = (p.Index - 1) % len(p.Data) // Circular pagination
		log.Printf("[INFO] Pagination %s data decremented", p.handerPrefix)
		p.UpdateMessage(s, i)
	}
}

func (p *PaginationView) PG_FirstBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("[INFO] Pagination %s data set to first", p.handerPrefix)
	p.Index = 0
	p.UpdateMessage(s, i)
}

func (p *PaginationView) PG_LastBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Printf("[INFO] Pagination %s data set to last index", p.handerPrefix)
	p.Index = len(p.Data) - 1
	p.UpdateMessage(s, i)
}

/*
Deletes the message if pressed
Note: this handler only works if the message is *not* ephemeral. Current impl. is
the message *is* ephemeral, so this handler has no shown corresponding button atm.
*/
func (p *PaginationView) PG_StopBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	/*
		Holy fucking shit
		this is how you delete your own fucking message
		why did it take so long to figure this out?
		fucks sake

		Do note: you can get the Channel ID of the Interaction with i.Message.ChannelID
		and it's own message id with i.Message.ID

	*/
	log.Printf("[INFO] Pagination %s stopped", p.handerPrefix)
	s.ChannelMessageDelete(i.Message.ChannelID, i.Message.ID)
	log.Printf("[INFO] Pagination %s removed from channel", p.handerPrefix)
}

/*
Gets rid of the buttons so theembed is permanent in the channel

# OR

Send a new message without the ephemeral flag and retain some buttons but not all (current impl.)
*/
func (p *PaginationView) PG_DoneBtnHandler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: p.CreateEmbed(),
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "View",
							Style: discordgo.LinkButton,
							URL:   p.EmbedLink,
						},
					},
				},
			},
		},
	})
	log.Printf("[INFO] Pagination %s stopped, embed remains in channel", p.handerPrefix)
}

/*
Add all the handlers for the buttons. Notably this is a direct copy-paste implementation from setup.go
Definitely not the best solution. Should probably change this later.
I am fairly confident (lol) that cleaner code here would not help the issue with concurrent paginations.
*/
func (p *PaginationView) PG_AddHandlers(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		// Attach component handlers, such as handlers for buttons
		case discordgo.InteractionMessageComponent:
			if h, ok := p.PageBtnHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	log.Printf("[INFO] Added pagination handlers to %s", p.handerPrefix)
}
