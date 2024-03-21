package music

import (
	"context"
	"discord-bot/handler"
	"discord-bot/services"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

type PauseHandler struct {
	Lavalink *services.Lavalink
}

func (h PauseHandler) Handle(session *discordgo.Session, event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := h.Lavalink.Client.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return handler.RespondWithContent(session, event.Interaction, "No player found")
	}

	if err := player.Update(context.TODO(), lavalink.WithPaused(!player.Paused())); err != nil {
		return handler.RespondWithContent(session, event.Interaction, fmt.Sprintf("Error while pausing: `%s`", err))
	}

	status := "playing"
	if player.Paused() {
		status = "paused"
	}

	return handler.RespondWithContent(session, event.Interaction, fmt.Sprintf("Player is now %s", status))
}
