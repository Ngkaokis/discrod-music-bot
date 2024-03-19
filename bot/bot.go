package bot

import (
	"context"
	"discord-bot/handler"
	"discord-bot/models"
	"discord-bot/util"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
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
	lavalinkClient := disgolink.New(snowflake.MustParse(session.State.User.ID))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	node, err := lavalinkClient.AddNode(ctx, disgolink.NodeConfig{
		Name:     "discord bot",
		Address:  config.LavalinkNodeAddress,
		Password: config.LavalinkPassword,
		Secure:   false,
	})
	if err != nil {
		log.Fatal(err)
	}
	version, err := node.Version(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("node version: %s", version)
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
		Queues:   *queueManager,
	}
	return bot, nil
}
