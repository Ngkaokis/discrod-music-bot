package util

import (
	"fmt"
	"strings"

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

func ParsePrefix(msg string ) (command string, query string, hasPrefix bool) {
	prefix := DiscordConfig.Prefix
	if strings.HasPrefix(msg, prefix) {
		command := strings.Replace(strings.Split(msg, " ")[0], prefix, "", 1)
		
		query := strings.TrimSpace(strings.Replace(msg, fmt.Sprintf("%s%s", prefix, command), "", 1))
		//lower case command
		command = strings.ToLower(command);
		return command ,query, true;

}
return "","",false
}

func GetGuildNameByID(bot *discordgo.Session, guildID string) string {
	guildName, ok := GuildNames[guildID]
	if !ok {
		guild, err := bot.Guild(guildID)
		if err != nil {
			// Failed to get the guild? Whatever, we'll just use the guild id
			GuildNames[guildID] = guildID
			return guildID
		}
		GuildNames[guildID] = guild.Name
		return guild.Name
	}
	return guildName
}
