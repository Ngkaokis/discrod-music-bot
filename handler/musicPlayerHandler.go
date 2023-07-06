package handler

import (
	"discord-bot/util"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)


func MusicPlayerHandler(s *discordgo.Session, m *discordgo.MessageCreate) {

	
	
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}
 
	command, query, hasPrefix := util.ParsePrefix(m.Content)

	if !hasPrefix {
		return
	}

	switch command {
		case "youtube" , "yt":
			fmt.Println(m.Content)
			fmt.Println(query)
			
		case "skip":
			//TODO SKIP implementation
		
		case "stop":
			//TODO STOP
		
		
		default: break;
	}

	// if m.Content == "Play" {
		
	// 	fmt.Println("working")
	
	// 	//find the voice channel of user
	// 	voiceChannel,err := util.SearchVoiceChannel(s , m)
	// 	if err != nil {
	// 		s.ChannelMessageSend(m.ChannelID, "OwO Please join the voice channel OwO")
	// 		fmt.Println("Can't find the voice channel",err)
	// 		return
	// 	}
	// 	fmt.Println("This is voice Channel:",voiceChannel )
		
		
		
	// 	vc ,err :=s.ChannelVoiceJoin(m.GuildID,voiceChannel,false,true)
	// 	if err != nil {
	// 		fmt.Println("Can't join the channel:",err)
	// 		return
	// 	}
		
		
	// 	dgvoice.PlayAudioFile(vc,"music.mp3",make(chan bool))
	// 	defer  vc.Disconnect()
	// }

	
	
}

func handleYoutubeCommand(s *discordgo.Session, m *discordgo.MessageCreate, query string) {
	voiceChannel,err := util.SearchVoiceChannel(s , m);
	if err!=nil{
		log.Printf("Can't find the voice channel",err);
	}
	fmt.Println(voiceChannel)
	//TODO
	
	
	}


