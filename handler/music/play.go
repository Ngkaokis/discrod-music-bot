package music

import (
	"context"
	"discord-bot/models"
	"discord-bot/services"
	"discord-bot/util"
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

type MusicPlayHandler struct {
	Lavalink     *services.Lavalink
	QueueManager *models.QueueManager
}

func (handler MusicPlayHandler) Handle(session *discordgo.Session, event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error {
	identifier := data.Options[0].StringValue()
	if !util.UrlPattern.MatchString(identifier) && !util.SearchPattern.MatchString(identifier) {
		identifier = lavalink.SearchTypeYouTube.Apply(identifier)
	}

	voiceState, err := session.State.VoiceState(event.GuildID, event.Member.User.ID)
	if err != nil {
		return session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Error while getting voice state: `%s`", err),
			},
		})
	}

	if err := session.InteractionRespond(event.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}); err != nil {
		return err
	}

	player := handler.Lavalink.Client.Player(snowflake.MustParse(event.GuildID))
	queue := handler.QueueManager.Get(event.GuildID)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var toPlay *lavalink.Track
	handler.Lavalink.Client.BestNode().LoadTracksHandler(ctx, identifier, disgolink.NewResultHandler(
		func(track lavalink.Track) {
			_, _ = session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Loading track: [`%s`](<%s>)", track.Info.Title, *track.Info.URI)),
			})
			if player.Track() == nil {
				toPlay = &track
			} else {
				queue.Add(track)
			}
		},
		func(playlist lavalink.Playlist) {
			_, _ = session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Loaded playlist: `%s` with `%d` tracks", playlist.Info.Name, len(playlist.Tracks))),
			})
			if player.Track() == nil {
				toPlay = &playlist.Tracks[0]
				queue.Add(playlist.Tracks[1:]...)
			} else {
				queue.Add(playlist.Tracks...)
			}
		},
		func(tracks []lavalink.Track) {
			_, _ = session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Loaded search result: [`%s`](<%s>)", tracks[0].Info.Title, *tracks[0].Info.URI)),
			})
			if player.Track() == nil {
				toPlay = &tracks[0]
			} else {
				queue.Add(tracks[0])
			}
		},
		func() {
			_, _ = session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Nothing found for: `%s`", identifier)),
			})
		},
		func(err error) {
			_, _ = session.InteractionResponseEdit(event.Interaction, &discordgo.WebhookEdit{
				Content: json.Ptr(fmt.Sprintf("Error while looking up query: `%s`", err)),
			})
		},
	))
	if toPlay == nil {
		return nil
	}

	if err := session.ChannelVoiceJoinManual(event.GuildID, voiceState.ChannelID, false, false); err != nil {
		return err
	}

	return player.Update(context.TODO(), lavalink.WithTrack(*toPlay))
}
