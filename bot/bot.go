package bot

import (
	"discord-bot/handler"
	"discord-bot/handler/music"
	"discord-bot/models"
	"discord-bot/util"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
)

type Bot struct {
	Session  *discordgo.Session
	Lavalink disgolink.Client
	Handlers map[string]handler.CommandHandler
	Queues   models.QueueManager
}

func New(session *discordgo.Session, config util.Config) (*Bot, error) {
	// Open a websocket connection to Discord and begin listening.
	err := session.Open()
	if err != nil {
		log.Fatal("error opening connection,", err)
		return nil, err
	}
	registerCommands(session)
	queueManager := models.NewQueueManger()
	lavalinkClient, err := registerLavalink(session, config)
  if err != nil {
    log.Fatal(err)
  }
	handlers := map[string]handler.CommandHandler{
		"play": &music.MusicPlayHandler{
			Lavalink:     lavalinkClient,
			QueueManager: queueManager,
		},
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
		Queues:   *queueManager,
	}
	return bot, nil
}
