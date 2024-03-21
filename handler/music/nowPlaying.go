package music

import (
	"discord-bot/handler"
	"discord-bot/services"
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

func formatPosition(position lavalink.Duration) string {
	if position == 0 {
		return "0:00"
	}
	return fmt.Sprintf("%d:%02d", position.Minutes(), position.SecondsPart())
}

type NowPlayingHandler struct {
	Lavalink *services.Lavalink
}

func (h NowPlayingHandler) Handle(session *discordgo.Session, event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	player := h.Lavalink.Client.ExistingPlayer(snowflake.MustParse(event.GuildID))
	if player == nil {
		return handler.RespondWithContent(session, event.Interaction, "No player found")
	}

	track := player.Track()
	if track == nil {
		return handler.RespondWithContent(session, event.Interaction, "No track found")
	}

	return handler.RespondWithContent(session, event.Interaction, fmt.Sprintf("Now playing: [`%s`](<%s>)\n\n %s / %s", track.Info.Title, *track.Info.URI, formatPosition(player.Position()), formatPosition(track.Info.Length)))
}
