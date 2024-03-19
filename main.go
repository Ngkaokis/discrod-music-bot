package main

import (
	"discord-bot/bot"
	"discord-bot/util"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	config util.Config
)

func init() {
	var err error
	config, err = util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not Load the Config File")
	}
}

func main() {

	// Create a new Discord session using the provided bot token.
	session, err := discordgo.New("Bot " + config.Token_String)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}

	b, err := bot.New(session, config)
	if err != nil {
		log.Fatal("Can not start the bot")
		return
	}

	// receiving message events and voice event.
	b.Session.Identify.Intents = discordgo.IntentsGuildMessages + discordgo.IntentsGuildVoiceStates

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	// Cleanly close down the Discord session.
	b.Session.Close()
}
