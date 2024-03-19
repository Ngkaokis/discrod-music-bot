package bot

import (
	"context"
	"discord-bot/util"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/snowflake/v2"
)

func registerLavalink (session *discordgo.Session, config util.Config) (disgolink.Client, error) {
	lavalinkClient := disgolink.New(snowflake.MustParse(session.State.User.ID))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	node, err := lavalinkClient.AddNode(ctx, disgolink.NodeConfig{
		Name:     "discord bot",
		Address:  config.LavalinkNodeAddress,
		Password: config.LavalinkPassword,
		Secure:   false,
	})
	if err != nil {
    return nil, err
	}
	version, err := node.Version(ctx)
	if err != nil {
    return nil, err
	}
	log.Printf("node version: %s", version)
  return lavalinkClient, nil
}
