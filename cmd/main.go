package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Potential-Central/NoExcusesBot"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	client, _ := NoExcusesBot.MakeBot(os.Getenv("DISCORD_TOKEN"))
	client.Start()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	client.Stop()
}
