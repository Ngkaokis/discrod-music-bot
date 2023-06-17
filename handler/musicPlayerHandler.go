package handler

import (
	"discord-bot/util"
	"fmt"

	"github.com/bwmarrin/discordgo"
)


func MusicPlayerHandler(s *discordgo.Session, m *discordgo.MessageCreate  ) {


	
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
 
fmt.Println(s.VoiceConnections)
	if m.Content == "Play" {
		
		fmt.Println("working")
	
		//find the voice channel of user
		voiceChannel,err := util.SearchVoiceChannel(s , m)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, "OwO Please join the voice channel OwO")
			fmt.Println("Can't find the voice channel",err)
			return
		}
		fmt.Println("This is voice Channel:",voiceChannel )
		
		
		
		vc ,err :=s.ChannelVoiceJoin(m.GuildID,voiceChannel,false,true)
		if err != nil {
			fmt.Println(err)
			return
		}
		
		
		// if err != nil {
		// 	time.Sleep(2 * time.Second)
		// 	vc, err = s.ChannelVoiceJoin("YOUR_GUILD_ID", "", false, false)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }

	vc.Close()
	// defer vc.Disconnect()
	}

	
}


