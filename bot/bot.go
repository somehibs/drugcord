package drugcord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/somehibs/tripapi/api"
	"regexp"
	"strings"
	"unicode"
)

type BotMain interface {
	Run() error
	Connect()
}

type Bot struct {
	ready   bool
	discord *discordgo.Session
	c       *BotConfig
}

var bot = Bot{ready: false, discord: nil}

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Bot is now READY.")
}

var mentionSyntax, _ = regexp.Compile("(\\<\\@[0-9]+\\>)+")

func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	m := mc.Message
	fmt.Printf("Message received. %+v\n", m)
	for _, v := range m.Mentions {
		if v.ID == bot.c.ID {
			fmt.Println("It's me! Check this for any commands.")
			bot.processMessage(m)
		}
	}
}

func StripContent(m *discordgo.Message) (content string) {
	content = m.Content
	for _, u := range m.Mentions {
		content = strings.NewReplacer("<@"+u.ID+">", "", "<@!"+u.ID+">", "").Replace(content)
	}
	return
}

func (b *Bot) processMessage(m *discordgo.Message) {
	// Every message contains the following entities
	// ChannelID, Timestamp, Content, EditedTimestamp, Tts, MentionEveryone, Attachments, Embeds, Mentions, Reactions, Type, ID
	// Read the content and return the drug formatted badly
	sm := StripContent(m)
	fmt.Println(sm)
	sm = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, sm)
	d := tripapi.GetDrug(sm)
	fmt.Printf("Found drug %+v\n", d)
}

func (b *Bot) addHandlers() {
	b.discord.AddHandler(onReady)
	b.discord.AddHandler(onMessageCreate)
}

func (b *Bot) Run() (e error) {
	// Get the configuration we're going to use.
	c := GetConf()
	b.c = &c
	fmt.Printf("e: %s p: %s t: %s\n", c.Email, c.Password, c.Token)
	fmt.Println("Initializing Discord session object.")
	b.discord, e = discordgo.New(c.Email, c.Password, c.Token)
	if e != nil {
		return Fatal(e, "Couldn't init Discord session obj.")
	}
	b.addHandlers()
	fmt.Println("Creating Discord session.")
	e = b.discord.Open()
	if e != nil {
		return Fatal(e, "Couldn't open a Discord session.")
	}
	return nil
}

func (b *Bot) Connect() {
}

func NewBot() *Bot {
	return &bot
}
