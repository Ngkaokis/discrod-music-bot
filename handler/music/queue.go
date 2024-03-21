package music

import (
	"discord-bot/handler"
	"discord-bot/models"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type QueueHandler struct {
	QueueMananger *models.QueueManager
}

func (h QueueHandler) Handle(session *discordgo.Session, event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := h.QueueMananger.Get(event.GuildID)
	if queue == nil {
		return handler.RespondWithContent(session, event.Interaction, "No player found")
	}

	if len(queue.Tracks) == 0 {
		return handler.RespondWithContent(session, event.Interaction, "No tracks in queue")
	}

	var tracks string
	for i, track := range queue.Tracks {
		tracks += fmt.Sprintf("%d. [`%s`](<%s>)\n", i+1, track.Info.Title, *track.Info.URI)
	}

	return handler.RespondWithContent(session, event.Interaction, fmt.Sprintf("Queue `%s`:\n%s", queue.Type, tracks))
}
