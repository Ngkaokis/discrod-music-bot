package util

import (
	"discord-bot/services"
	"discord-bot/types"
	"sync"
)

var  (

	Guilds 		= make(map[string]*types.ActiveGuild)
	GuildsMutex = sync.RWMutex{}
	GuildNames   = make(map[string]string)

	YoutubeService = services.NewService(1000)

)
