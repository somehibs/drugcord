package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/somehibs/drugcord/bot"
)

const testBotEnabled = true

func main() {
	// Load bot instance
	var testBot *drugcord.Bot = nil
	if testBotEnabled {
		//fmt.Println(strings.Join(drugcord.DiscordFormatter{}.FormatTableFields(tripapi.GetDrug("heroin")), "\n"))
		//panic("Okay.")
		testBot = drugcord.NewBot(drugcord.GetConfByName("./testconfig.json"))
		testBot.Connect()
	}

	bot := drugcord.NewBot(drugcord.GetConf())
	err := bot.Connect()
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	bot.RouteCommands()
	_, e := testBot.Discord.ChannelMessageSend("334458647788912640", "!drug heroin")
	if e != nil {
		fmt.Printf("couldn't send using testbot %s\n", e)
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, os.Kill)
	<-sigChan
	//	fmt.Println(tripapi.GetDrug("molly"))
}
