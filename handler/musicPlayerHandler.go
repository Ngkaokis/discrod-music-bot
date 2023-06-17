package handler

import (
	"discord-bot/types"
	"discord-bot/util"
	"discord-bot/workers"
	"os"

	"fmt"

	"github.com/bwmarrin/discordgo"
)


func MusicPlayerHandler(s *discordgo.Session, m *discordgo.MessageCreate  ) {

	
	
	// Ignore all messages created by the bot itself
	if message.Author.ID == bot.State.User.ID {
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
			fmt.Println("Can't join the channel:",err)
			return
		}
		
		
		// if err != nil {
		// 	time.Sleep(2 * time.Second)
		// 	vc, err = s.ChannelVoiceJoin("YOUR_GUILD_ID", "953675829668347964", false, false)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }
		
		dgvoice.PlayAudioFile(vc,"music.mp3",make(chan bool))
		defer  vc.Disconnect()
	}

	

	
}



		// Append encoded pcm data to the buffer.
		buffer = append(buffer, InBuf)
	}
}