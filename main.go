package main

import (
	"os"
	"os/signal"
	"syscall"

	bot "github.com/Potential-Central/NoExcusesBot/bot"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	client, _ := bot.MakeBot(os.Getenv("DISCORD_TOKEN"))
	client.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	client.Stop()
}
