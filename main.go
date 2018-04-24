package main

import (
	"fmt"
	"github.com/somehibs/drugcord/bot"
	_ "github.com/somehibs/tripapi/api"
	"os"
	"os/signal"
)

const testBotEnabled = true

func main() {
	// Load bot instance
	var testBot *drugcord.Bot = nil
	if testBotEnabled {
		testBot = drugcord.NewBot(drugcord.GetConfByName("./testconfig.json"))
		testBot.Connect()
	}

	bot := drugcord.NewBot(drugcord.GetConf())
	err := bot.Connect()
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		_, e := testBot.Discord.ChannelMessageSend("438143058807357470", "!drug mdma")
		if e != nil {
			fmt.Printf("%s\n", e)
		}
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, os.Kill)
		<-sigChan
	}
	//	fmt.Println(tripapi.GetDrug("molly"))
}
