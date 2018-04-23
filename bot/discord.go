package drugcord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	_ "github.com/somehibs/tripapi/api"
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
	cmd     CommandRouter
}

var bot = Bot{ready: false, discord: nil}

const cmdChar = '!'

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	fmt.Println("Bot is now READY.")
}

func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	m := mc.Message
	//fmt.Printf("Message received %+v\n", m)
	if m.Content[0] == cmdChar {
		bot.processMessage(m)
		return
	}
	for _, v := range m.Mentions {
		if v.ID == bot.c.ID {
			m.Content = StripMentions(m)
			bot.processMessage(m)
		}
	}
}

func StripMentions(m *discordgo.Message) (content string) {
	content = m.Content
	for _, u := range m.Mentions {
		content = strings.NewReplacer("<@"+u.ID+">", "", "<@!"+u.ID+">", "").Replace(content)
	}
	if content[0] == ' ' {
		content = content[1:]
	}
	return
}

func StripSpace(s string) string {
	return strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, s)
}

func (b *Bot) processMessage(m *discordgo.Message) {
	// Every message contains the following entities
	// ChannelID, Timestamp, Content, EditedTimestamp, Tts, MentionEveryone, Attachments, Embeds, Mentions, Reactions, Type, ID
	if m.Content[0] == cmdChar {
		m.Content = m.Content[1:]
	}
	fmt.Printf("Beginning to process: %s\n", m.Content)
	input := MessageInput{OriginalMessage: m, Content: m.Content}
	go b.cmd.HandleMessage(&input)
}

func (b *Bot) addHandlers() {
	b.discord.AddHandler(onReady)
	b.discord.AddHandler(onMessageCreate)
}

func (b Bot) Send(responses CommandResponse) {
}

func (b Bot) SendAll(responses []CommandResponse) {
	for response := range responses {
		fmt.Println("Response %s", response)
	}
}

func (b *Bot) Run() (e error) {
	//fmt.Printf("e: %s p: %s t: %s\n", c.Email, c.Password, c.Token)
	fmt.Println("Initializing Discord session object.")
	b.discord, e = discordgo.New(b.c.Email, b.c.Password, b.c.Token)
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

func NewBot() *Bot {
	// Get the configuration we're going to use, init other things.
	c := GetConf()
	bot.c = &c
	bot.cmd = CommandRouter{Handler: bot}
	bot.cmd.RegisterCommands(DrugCommands)
	return &bot
}
