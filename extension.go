package NoExcusesBot

import (
	"github.com/bwmarrin/discordgo"
)

type Extension interface {
	CommandList() []*Command
	Name() string
	Help() string
}

type Command struct {
	Name          string
	HasPermission func(s *discordgo.Session, m *discordgo.MessageCreate) bool
	Execute       func(s *discordgo.Session, m *discordgo.MessageCreate)
}
