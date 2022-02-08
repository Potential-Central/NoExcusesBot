package main

import (
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	logger       *log.Logger
	database     *sql.DB
	client       *discordgo.Session
	tasks        map[int]Task  //TaskID  -> Task
	guilds       map[int]Guild //GuildID -> Guild
	pendingTasks map[int]Task  //UserID  -> Task
	help         map[string]discordgo.MessageEmbed
)

func init() {
	logger = log.Default()
	err := godotenv.Load()
	if err != nil {
		logger.Fatal("[SETUP] Error loading .env file")
	}
	createDatabase()
	getGuilds()
	getTasks()
	GetHelp()
	CompileRegex()
}

func main() {
	var err error
	defer database.Close()

	//Making task handler
	tz, _ := time.LoadLocation("UTC")
	scheduler := gocron.NewScheduler(tz)
	scheduler.Every("1m").SingletonMode().Do(CheckTasks)

	//Making discord client
	client, err = discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
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

	//Starting tasks
	scheduler.StartAsync()

	logger.Println("[SETUP] Now running.  Press CTRL-C to exit.")
	//Gracefully close from console
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	logger.Println("[SETUP] Shutting down...")
	client.Close()
	scheduler.Stop()
}
