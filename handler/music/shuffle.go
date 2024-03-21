package music

import (
	"discord-bot/handler"
	"discord-bot/models"

	"github.com/bwmarrin/discordgo"
)

type ShuffleHandler struct {
	QueueMananger *models.QueueManager
}

func (h ShuffleHandler) Handle(session *discordgo.Session, event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	queue := h.QueueMananger.Get(event.GuildID)
	if queue == nil {
		return handler.RespondWithContent(session, event.Interaction, "No player found")
	}

	queue.Shuffle()
	return handler.RespondWithContent(session, event.Interaction, "Queue shuffled")
}
