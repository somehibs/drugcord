package main

import (
	"fmt"
	"github.com/somehibs/drugcord/bot"
	_ "github.com/somehibs/tripapi/api"
	"os"
	"os/signal"
)

const testBotEnabled = false

func main() {
	// Load bot instance
	if testBotEnabled {
		testBot := drugcord.NewBot(drugcord.BotConfig{Token: "Bot UNKNOWN"})
		testBot.Connect()
	}

	bot := drugcord.NewBot(drugcord.GetConf())
	err := bot.Connect()
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, os.Kill)
		<-sigChan
	}
	//	fmt.Println(tripapi.GetDrug("molly"))
}
