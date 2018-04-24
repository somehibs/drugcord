package drugcord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"unicode"
)

type Bot struct {
	ready   bool
	discord *discordgo.Session
	Discord *discordgo.Session
	user    *discordgo.User
	c       *BotConfig
	cmd     CommandRouter
}

const cmdChar = '!'

var bots = map[*discordgo.Session]*Bot{}
var firstbot *Bot = nil

func onReady(s *discordgo.Session, event *discordgo.Ready) {
	b := bots[s]
	if b == nil {

	}
	fmt.Printf("Bot ready: %s\n", b.c.ID)
}

func botFromSession(s *discordgo.Session) *Bot {
	return bots[s]
}

func onMessageCreate(s *discordgo.Session, mc *discordgo.MessageCreate) {
	bot := botFromSession(s)
	if bot == nil {
		fmt.Printf("Could not find bot for session %+v\n", s)
		return
	}
	if bot.user != nil && bot.user.ID == mc.Author.ID {
		return
	}
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

func (b Bot) processMessage(m *discordgo.Message) {
	// Every message contains the following entities
	// ChannelID, Timestamp, Content, EditedTimestamp, Tts, MentionEveryone, Attachments, Embeds, Mentions, Reactions, Type, ID
	if m.Content[0] == cmdChar {
		m.Content = m.Content[1:]
	}
	fmt.Printf("Command received: %s\n", m.Content)
	input := MessageInput{Original: m, Content: m.Content}
	go b.cmd.HandleMessage(b, &input)
}

func (b Bot) addHandlers() {
	b.discord.AddHandler(onReady)
	b.discord.AddHandler(onMessageCreate)
}

func (b Bot) Send(response CommandResponse) {
	// Find out who to send it to
	m := response.Input.Original.(*discordgo.Message)
	message := strings.Join(response.Reply, "")
	//msg := discordgo.MessageSend{Content: message}
	switch response.Target {
	case TargetAdminChannel:
	case TargetOtherChannel:
	case TargetNone:
	default:
		fmt.Printf("Target: %s will not receive %s", message)
	case TargetRequestor:
	case TargetSameChannel:
		fmt.Printf("Send: %s %s\n", m.ChannelID, message)
		v, e := b.discord.ChannelMessageSend(m.ChannelID, message)
		if e != nil {
			fmt.Printf("Error %+v\n\n", v, e)
		}
		//fmt.Printf("Sent: %+v %+v\n\n", v, e)
	}
	//fmt.Printf("Send %s to %s\n", message, response.Target)
}

func (b Bot) SendAll(responses []CommandResponse) {
	for _, response := range responses {
		b.Send(response)
	}
}

func (b *Bot) Connect() (e error) {
	fmt.Println("Initializing Discord session object.")
	b.discord, e = discordgo.New(b.c.Token)
	b.Discord = b.discord
	if e != nil {
		return Fatal(e, "Couldn't init Discord session obj.")
	}

	b.addHandlers()
	bots[b.discord] = b

	fmt.Println("Creating Discord session.")
	e = b.discord.Open()
	if e != nil {
		delete(bots, b.discord)
		return Fatal(e, "Couldn't open a Discord session.")
	}

	user, e := b.discord.User("@me")
	if e != nil {
		return Fatal(e, "Could not fetch user for session.")
	}
	b.user = user

	// Handle some commands with a router
	b.cmd = CommandRouter{}
	b.cmd.RegisterCommands(DrugCommands)

	return nil
}

func NewBot(c BotConfig) *Bot {
	// Get the configuration we're going to use, init other things.
	bot := Bot{ready: false, discord: nil}
	bot.c = &c

	return &bot
}
