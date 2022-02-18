package main

import (
	"os"
	"os/signal"
	"syscall"

	bot "github.com/Potential-Central/NoExcusesBot"
	exts "github.com/Potential-Central/NoExcusesBot/exts"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	client, _ := bot.MakeBot(os.Getenv("DISCORD_TOKEN"))
	exts.MakeChannelsExt(client)

	client.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	client.Stop()
}
