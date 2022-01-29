package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	logger   *log.Logger
	database *sql.DB
	client   *discordgo.Session
	tasks    map[int]Task
	guilds   map[int]Guild
)

func init() {
	logger = log.Default()
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("[SETUP] Error loading .env file")
	}
	createDatabase()
	getGuilds()
}

func main() {
	defer database.Close()
	client, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		logger.Fatal("[SETUP] Error opening connection,", err)
	}
	client.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages

	//Commands handler
	client.AddHandler(messageHandler)

	//Guild join handler
	client.AddHandler(guildJoinHandler)

	err = client.Open()
	if err != nil {
		logger.Fatal("[SETUP] Error opening connection,", err)
	}

	logger.Println("[SETUP] Now running.  Press CTRL-C to exit.")
	//Gracefully close from console
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	logger.Println("[SETUP] Shutting down...")
	client.Close()
}
