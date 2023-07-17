package util

import (
	"discord-bot/models"
	"discord-bot/services"
	"sync"
)

var  (

	Guilds 		= make(map[string]*models.ActiveGuild)
	GuildsMutex = sync.RWMutex{}
	GuildNames   = make(map[string]string)

	YoutubeService = services.NewService(1000)

)
