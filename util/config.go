package util

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Token_String string `mapstructure:"DISCORD_BOT_TOKEN"`
	Prefix string `mapstructure:"PREFIX"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("bot")
	viper.SetConfigType("env")

	
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return 
}

func Get() (Config) {
	cfg,err := LoadConfig(".")
	if(err != nil){
		log.Fatal(err)
	}
	return cfg
}