package main

import (
	"fmt"
	"github.com/somehibs/drugcord/bot"
	_ "github.com/somehibs/tripapi/api"
	"os"
	"os/signal"
)

func main() {
	// Load bot instance
	err := drugcord.NewBot().Run()
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, os.Kill)
		<-sigChan
	}
	//	fmt.Println(tripapi.GetDrug("molly"))
}
