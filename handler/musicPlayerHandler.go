package handler

import (
	"discord-bot/types"
	"discord-bot/util"
	"discord-bot/workers"
	"os"

	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)


func MusicPlayerHandler(bot *discordgo.Session, message *discordgo.MessageCreate) {

	
	
	// Ignore all messages created by the bot itself
	if message.Author.ID == bot.State.User.ID {
		return
	}
 
	command, query, hasPrefix := util.ParsePrefix(message.Content)
	if !hasPrefix {
		return
	} 
		
	util.GuildsMutex.Lock()
	activeGuild := util.Guilds[message.GuildID]
	util.GuildsMutex.Unlock()
	
	switch command {
		case "youtube" , "yt":
			handleYoutubeCommand(bot, message, activeGuild, query)
			
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

func handleYoutubeCommand(bot *discordgo.Session,message *discordgo.MessageCreate,
	activeGuild *types.ActiveGuild,  query string) {
	
	//* handle the queue is full
	if activeGuild != nil{
		if activeGuild.IsMediaQueueFull(){
			_, _ = bot.ChannelMessageSend(message.ChannelID, "The queue is full!")
			return
		}
	} else {
		//* If first time to play music, add the guild
		activeGuild = types.NewActiveGuild(util.GetGuildNameByID(bot, message.GuildID))
		util.GuildsMutex.Lock()
		util.Guilds[message.GuildID] = activeGuild
		util.GuildsMutex.Unlock()
	}
	
	voiceChannelId,err := util.SearchVoiceChannel(bot , message);
	if err!=nil{
		log.Printf("[%s] Can't find the voice channel",err);
	} else {
		log.Printf("Success");
	}
	fmt.Println(voiceChannelId)

	log.Printf("[%s] Searching for \"%s\"", activeGuild.Name, query)

	//* Search music
	media, err := util.YoutubeService.SearchYTAndDownload(query)
	if err != nil {
		log.Printf("[%s] Unable to find video for query \"%s\": %s", activeGuild.Name, query, err.Error())
		return
	}

	//* Add song to guild queue
	createNewWorker := false
	if !activeGuild.IsStreaming() {
		log.Printf("[%s] Preparing for streaming", activeGuild.Name)
		activeGuild.PrepareForStreaming(20)
		// If the channel was nil, it means that there was no worker
		createNewWorker = true
	}
	activeGuild.EnqueueMedia(media)

	if createNewWorker {
		log.Printf("[%s] Starting worker", activeGuild.Name)
		go func() {
			err = workers.Worker(bot, activeGuild, message.GuildID, voiceChannelId)
			if err != nil {
				log.Printf("[%s] Failed to start worker: %s", activeGuild.Name, err.Error())
				_, _ = bot.ChannelMessageSend(message.ChannelID, fmt.Sprintf("❌ Unable to start voice worker: %s", err.Error()))
				_ = os.Remove(media.FilePath)
				return
			}
		}()
	}

	}


