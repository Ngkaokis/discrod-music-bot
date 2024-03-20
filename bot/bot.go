package bot

import (
	"context"
	"discord-bot/handler"
	"discord-bot/handler/music"
	"discord-bot/models"
	"discord-bot/services"
	"discord-bot/util"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"
)

type Bot struct {
	Session      *discordgo.Session
	Lavalink     *services.Lavalink
	Handlers     map[string]handler.CommandHandler
	QueueManager *models.QueueManager
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
	lavalink, err := services.NewLavaLinkService(session, queueManager, config)
	if err != nil {
		log.Fatal(err)
	}
	handlers := map[string]handler.CommandHandler{
		"play": &music.MusicPlayHandler{
			Lavalink:     lavalink,
			QueueManager: queueManager,
		},
	}
	bot := &Bot{
		Session:      session,
		Lavalink:     lavalink,
		Handlers:     handlers,
		QueueManager: queueManager,
	}
	bot.Session.AddHandler(bot.onApplicationCommand)
	bot.Session.AddHandler(bot.onVoiceStateUpdate)
	bot.Session.AddHandler(bot.onVoiceServerUpdate)
	return bot, nil
}

func (b *Bot) onApplicationCommand(session *discordgo.Session, event *discordgo.InteractionCreate) {
	data := event.ApplicationCommandData()

	handler, ok := b.Handlers[data.Name]
	if !ok {
		log.Printf("unknown command: %s", data.Name)
		return
	}
	if err := handler.Handle(session, event, data); err != nil {
		log.Printf("error handling command: %s", err)
	}
}

func (b *Bot) onVoiceStateUpdate(session *discordgo.Session, event *discordgo.VoiceStateUpdate) {
	if event.UserID != session.State.User.ID {
		return
	}

	var channelID *snowflake.ID
	if event.ChannelID != "" {
		id := snowflake.MustParse(event.ChannelID)
		channelID = &id
	}
	b.Lavalink.Client.OnVoiceStateUpdate(context.TODO(), snowflake.MustParse(event.GuildID), channelID, event.SessionID)
	if event.ChannelID == "" {
		b.QueueManager.Delete(event.GuildID)
	}
}

func (b *Bot) onVoiceServerUpdate(_ *discordgo.Session, event *discordgo.VoiceServerUpdate) {
	b.Lavalink.Client.OnVoiceServerUpdate(context.TODO(), snowflake.MustParse(event.GuildID), event.Token, event.Endpoint)
}
