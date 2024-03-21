package music

import (
	"discord-bot/handler"
	"discord-bot/services"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/snowflake/v2"
)

type StopHandler struct {
	Lavalink *services.Lavalink
}

func (h StopHandler) Handle(session *discordgo.Session, event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := h.Lavalink.Client.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return handler.RespondWithContent(session, event.Interaction, "No player found")
	}

	if err := session.ChannelVoiceJoinManual(event.GuildID, "", false, false); err != nil {
		return handler.RespondWithContent(session, event.Interaction, fmt.Sprintf("Error while disconnecting: `%s`", err))
	}

	return handler.RespondWithContent(session, event.Interaction, "Player stopped")
}
