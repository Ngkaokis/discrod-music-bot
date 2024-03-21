package handler

import "github.com/bwmarrin/discordgo"

func RespondWithContent(session *discordgo.Session, interaction *discordgo.Interaction, content string) error {
	return session.InteractionRespond(interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: content,
		},
	})
}
