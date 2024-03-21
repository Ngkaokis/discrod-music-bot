package music

import (
	"discord-bot/handler"
	"discord-bot/models"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type QueueTypeHandler struct {
	QueueMananger *models.QueueManager
}

func (h QueueTypeHandler) Handle(session *discordgo.Session, event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := h.QueueMananger.Get(event.GuildID)
	if queue == nil {
		return handler.RespondWithContent(session, event.Interaction, "No player found")
	}

	queue.Type = models.QueueType(data.Options[0].Value.(string))
	return handler.RespondWithContent(session, event.Interaction, fmt.Sprintf("Queue type set to `%s`", queue.Type))
}
