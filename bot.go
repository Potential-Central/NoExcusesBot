package main

import (
    "log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var (
	logger *log.Logger
	client *discordgo.Session
)

func init() {
	logger = log.Default()
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("[SETUP] Error loading .env file")
  	}
}

func main() {
	client, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		logger.Fatal("[SETUP] Error opening connection,", err)
	}
	client.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages
}