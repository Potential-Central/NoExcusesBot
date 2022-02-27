package NoExcusesBot

import (
	"strconv"

	"github.com/Potential-Central/NoExcusesBot"
	"github.com/bwmarrin/discordgo"
)

type TasksExt struct {
	Bot      *NoExcusesBot.Bot
	Commands []*NoExcusesBot.Command
}

func MakeTasksExt(bot *NoExcusesBot.Bot) {
	ret := &TasksExt{bot, make([]*NoExcusesBot.Command, 0)}
	//Registering extension commands

	bot.Exts = append(bot.Exts, ret)
	bot.Logger.Printf("[TASKS] Extension loaded")
}

func (ext *TasksExt) Name() string {
	return "TasksExt"
}

func (ext *TasksExt) Help() string {
	return "TasksExt Help"
}

func (ext *TasksExt) CommandList() []*NoExcusesBot.Command {
	return ext.Commands
}

func (ext *TasksExt) userPerms(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	if guild, ok := ext.Bot.Guilds[m.GuildID]; ok {
		intC, _ := strconv.Atoi(m.ChannelID)
		if intC == guild.UserChannel || intC == guild.AdminChannel {
			return true
		}
	}
	return false
}
