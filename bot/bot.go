package bot

import (
	"discord-bot/handler"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

type Bot struct {
	Session  *discordgo.Session
	Lavalink disgolink.Client
	Handlers map[string]handler.CommandHandler
}

func New(session *discordgo.Session) (*Bot, error) {
	// Open a websocket connection to Discord and begin listening.
  err := session.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
    return nil, err
	}
	registerCommands(session)
	lavalinkClient := disgolink.New(snowflake.MustParse(session.State.User.ID))
	handlers := map[string]handler.CommandHandler{
		"play": nil,
	}
	session.AddHandler(func(session *discordgo.Session, event *discordgo.InteractionCreate) {
		data := event.ApplicationCommandData()

		handler, ok := handlers[data.Name]
		if !ok {
			log.Printf("unknown command: %s", data.Name)
			return
		}
    log.Printf("Received command %s", data.Name)
		if err := handler.Handle(session, event, data); err != nil {
			log.Printf("error handling command: %+v", err)
		}
	})
	bot := &Bot{
		Session:  session,
		Lavalink: lavalinkClient,
		Handlers: handlers,
	}
	return bot, nil
}
