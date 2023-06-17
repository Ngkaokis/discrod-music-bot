package main

import (
	"discord-bot/handler"
	"discord-bot/util"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (config util.Config)

func init() {
	var err error
	config ,err  = util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not Load the Config File")
	}
}


func main() {

	// Create a new Discord session using the provided bot token.

	dg, err := discordgo.New("Bot " + config.Token_String)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(handler.MessageCreateHandler)

	dg.AddHandler(handler.MusicPlayerHandler)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.IntentsGuildMessages + discordgo.IntentsGuildVoiceStates

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

