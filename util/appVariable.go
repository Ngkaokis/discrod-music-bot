package util

import (
	"discord-bot/models"
	"sync"
)

var (
	Guilds      = make(map[string]*models.ActiveGuild)
	GuildsMutex = sync.RWMutex{}
	GuildNames  = make(map[string]string)
)
