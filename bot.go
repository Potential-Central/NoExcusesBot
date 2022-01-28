package main

import (
	"database/sql"
    "log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	logger *log.Logger
	database *sql.DB
	client *discordgo.Session
)

func init() {
	logger = log.Default()
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("[SETUP] Error loading .env file")
  	}
	database, _ = sql.Open("sqlite3", "tasks.db")
	dat, err := os.ReadFile("queries/createTables.sql")
	if err != nil {
		logger.Fatal("[SETUP] Error creating SQL tables")
	}
	database.Exec(string(dat))
}

func main() {
	defer database.Close()
	client, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		logger.Fatal("[SETUP] Error opening connection,", err)
	}
	client.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages
}