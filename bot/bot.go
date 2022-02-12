package NoExcusesBot

import (
	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	Perfix  string
	Token   string
	Session *discordgo.Session
}
