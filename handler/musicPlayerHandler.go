package handler

import (
	"discord-bot/models"
	"discord-bot/util"
	"discord-bot/workers"
	"fmt"
	"os"

	"log"

	"github.com/bwmarrin/discordgo"
)


func MusicPlayerHandler(bot *discordgo.Session, message *discordgo.MessageCreate) {

	
	if message.Content == "Play" {
		
		fmt.Println("working")
	
		//find the voice channel of user
		voiceChannel,err := util.SearchVoiceChannel(bot , message)
		if err != nil {
			bot.ChannelMessageSend(message.ChannelID, "OwO Please join the voice channel OwO")
			fmt.Println("Can't find the voice channel",err)
			return
		}
			
		fmt.Println("This is voice Channel:",voiceChannel )
		
		
		
		vc ,err :=bot.ChannelVoiceJoin(message.GuildID,voiceChannel,false,true)
		if err != nil {
			fmt.Println("Can't join the channel:",err)
			return
		}
		
		
		// dgvoice.PlayAudioFile(vc,"",make(chan bool))
		// defer  vc.Disconnect()
	}
	
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


	
	
}

func handleYoutubeCommand(bot *discordgo.Session,message *discordgo.MessageCreate,
	activeGuild *models.ActiveGuild,  query string) {
	
	//* handle the queue is full
	if activeGuild != nil{
		if activeGuild.IsMediaQueueFull(){
			_, _ = bot.ChannelMessageSend(message.ChannelID, "The queue is full!")
			return
		}
	} else {
		//* If first time to play music, add the guild
		activeGuild = models.NewActiveGuild(util.GetGuildNameByID(bot, message.GuildID))
		util.GuildsMutex.Lock()
		util.Guilds[message.GuildID] = activeGuild
		util.GuildsMutex.Unlock()
	}
	
	voiceChannelId,err := util.SearchVoiceChannel(bot , message);
	if err!=nil{
		log.Printf("[%s] Can't find the voice channel",err);
		bot.ChannelMessageSend(message.ChannelID, err.Error())
		return
	} else {
		log.Printf("Success");
	}
	log.Println(voiceChannelId)

	log.Printf("[%s] Searching for \"%s\"", activeGuild.Name, query)

	//* Search music
	media, err := util.YoutubeService.SearchYTAndDownload(query)
	if err != nil {
		log.Printf("[%s] Unable to find video for query \"%s\": %s", activeGuild.Name, query, err.Error())
		return
	}

	createNewWorker := false
	if !activeGuild.IsStreaming() {
		log.Printf("[%s] Preparing for streaming", activeGuild.Name)
		activeGuild.PrepareForStreaming(20)
		// If the channel was nil, it means that there was no worker
		createNewWorker = true
	}

	//* Add song to guild queue
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


