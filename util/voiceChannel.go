package util

import (
	"github.com/bwmarrin/discordgo"
)

func SearchVoiceChannel(s *discordgo.Session,m *discordgo.MessageCreate)(voiceChannelID string, err error){
	guild, err:= s.State.Guild(m.GuildID)
	if err != nil {
		return "", err
	}
		for _, vs := range guild.VoiceStates{
			if vs.UserID == m.Author.ID {
				// fmt.Println("vc:",vs.ChannelID)
				return vs.ChannelID,nil
			}
		
	}
	return "",nil
}