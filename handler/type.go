package handler

import (
	"github.com/bwmarrin/discordgo"
)

type CommandHandler interface {
	Handle(session *discordgo.Session, event *discordgo.InteractionCreate, data discordgo.ApplicationCommandInteractionData) error
}
