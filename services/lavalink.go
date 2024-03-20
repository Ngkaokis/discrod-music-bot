package services

import (
	"context"
	"discord-bot/models"
	"discord-bot/util"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type Lavalink struct {
	Client       disgolink.Client
	QueueManager *models.QueueManager
}

func NewLavaLinkService(session *discordgo.Session, queueManager *models.QueueManager, config util.Config) (*Lavalink, error) {
	lavalink := &Lavalink{
		QueueManager: queueManager,
	}
	lavalinkClient := disgolink.New(snowflake.MustParse(session.State.User.ID),
		disgolink.WithListenerFunc(lavalink.onPlayerPause),
		disgolink.WithListenerFunc(lavalink.onPlayerResume),
		disgolink.WithListenerFunc(lavalink.onTrackStart),
		disgolink.WithListenerFunc(lavalink.onTrackEnd),
		disgolink.WithListenerFunc(lavalink.onTrackException),
		disgolink.WithListenerFunc(lavalink.onTrackStuck),
		disgolink.WithListenerFunc(lavalink.onWebSocketClosed),
	)
	lavalink.Client = lavalinkClient

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	node, err := lavalinkClient.AddNode(ctx, disgolink.NodeConfig{
		Name:     "discord bot",
		Address:  config.LavalinkNodeAddress,
		Password: config.LavalinkPassword,
		Secure:   false,
	})
	if err != nil {
		return nil, err
	}
	version, err := node.Version(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("node version: %s", version)
	return &Lavalink{
		Client:       lavalinkClient,
		QueueManager: queueManager,
	}, nil
}

func (l *Lavalink) onPlayerPause(player disgolink.Player, event lavalink.PlayerPauseEvent) {
	log.Printf("onPlayerPause: %v\n", event)
}

func (l *Lavalink) onPlayerResume(player disgolink.Player, event lavalink.PlayerResumeEvent) {
	log.Printf("onPlayerResume: %v\n", event)
}

func (l *Lavalink) onTrackStart(player disgolink.Player, event lavalink.TrackStartEvent) {
	log.Printf("onTrackStart: %v\n", event)
}

func (l *Lavalink) onTrackEnd(player disgolink.Player, event lavalink.TrackEndEvent) {
	log.Printf("onTrackEnd: %v\n", event)

	if !event.Reason.MayStartNext() {
		return
	}

	queue := l.QueueManager.Get(event.GuildID().String())
	var (
		nextTrack lavalink.Track
		ok        bool
	)
	switch queue.Type {
	case models.QueueTypeNormal:
		nextTrack, ok = queue.Next()

	case models.QueueTypeRepeatTrack:
		nextTrack = event.Track

	case models.QueueTypeRepeatQueue:
		queue.Add(event.Track)
		nextTrack, ok = queue.Next()
	}

	if !ok {
		return
	}
	if err := player.Update(context.TODO(), lavalink.WithTrack(nextTrack)); err != nil {
		log.Printf("Failed to play next track: %s\n", err)
	}
}

func (l *Lavalink) onTrackException(player disgolink.Player, event lavalink.TrackExceptionEvent) {
	log.Printf("onTrackException: %v\n", event)
}

func (l *Lavalink) onTrackStuck(player disgolink.Player, event lavalink.TrackStuckEvent) {
	log.Printf("onTrackStuck: %v\n", event)
}

func (l *Lavalink) onWebSocketClosed(player disgolink.Player, event lavalink.WebSocketClosedEvent) {
	log.Printf("onWebSocketClosed: %v\n", event)
}
