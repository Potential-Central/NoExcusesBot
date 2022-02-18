package NoExcusesBot

import (
	"strconv"

	"github.com/Potential-Central/NoExcusesBot"
	"github.com/bwmarrin/discordgo"
)

type ChanExt struct {
	Bot      *NoExcusesBot.Bot
	Commands []*NoExcusesBot.Command
}

func MakeChannelsExt(bot *NoExcusesBot.Bot) {
	ret := &ChanExt{bot, make([]*NoExcusesBot.Command, 0)}
	//Registering extension commands
	ret.Commands = append(ret.Commands, &NoExcusesBot.Command{
		Name: "userchan", HasPermission: ret.adminPerms, Execute: ret.setUserChannel,
	})
	ret.Commands = append(ret.Commands, &NoExcusesBot.Command{
		Name: "adminchan", HasPermission: ret.adminPerms, Execute: ret.setAdminChannel,
	})
	bot.Exts = append(bot.Exts, ret)
	bot.Logger.Printf("[CHANS] Extension loaded")
}

func (ext *ChanExt) Name() string {
	return "ChannelsExt"
}

func (ext *ChanExt) Help() string {
	return "ChannelsExt Help"
}

func (ext *ChanExt) CommandList() []*NoExcusesBot.Command {
	return ext.Commands
}

//Checks if user has admin permissions
func (ext *ChanExt) adminPerms(s *discordgo.Session, m *discordgo.MessageCreate) bool {
	//Getting user permissions
	perm, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		return false
	}
	//Checking if usert is guild admin
	if perm&discordgo.PermissionAdministrator != 0 {
		return true
	}
	intCh, _ := strconv.Atoi(m.ChannelID)
	if guild, ok := ext.Bot.Guilds[m.GuildID]; ok {
		//Check for admin channel if one exists
		if intCh == guild.AdminChannel || guild.AdminChannel == 0 {
			//Check if user has admin role
			for _, role := range m.Member.Roles {
				intR, _ := strconv.Atoi(role)
				if intR == guild.AdminRole {
					return true
				}
			}
		}
	}
	return false
}

//Setting a channel where reminder get posted, and users interact with the bot
func (ext *ChanExt) setUserChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
	intG, _ := strconv.Atoi(m.GuildID)
	intCh, _ := strconv.Atoi(m.ChannelID)
	if guild, ok := ext.Bot.Guilds[m.GuildID]; ok {
		//If guild already exists, just update the channel.
		guild.UserChannel = intCh
		NoExcusesBot.WriteObject(ext.Bot.Database, guild)
		s.ChannelMessageSend(m.ChannelID, "This channel is now the reminder channel!")
		ext.Bot.Logger.Printf("[CHANS] Updated user channel in guild %v to %v", intG, intCh)
	} else {
		//If guild doesn't exist, create new guild
		guild = &NoExcusesBot.Guild{
			Id:          intG,
			UserChannel: intCh,
		}
		ext.Bot.Guilds[m.GuildID] = guild
		NoExcusesBot.WriteObject(ext.Bot.Database, guild)
		s.ChannelMessageSend(m.ChannelID, "This channel is now the reminder channel!")
		ext.Bot.Logger.Printf("[CHANS] Created user channel in new guild %v to %v", intG, intCh)
	}
}

//Setting a channel where administrative commands are used
func (ext *ChanExt) setAdminChannel(s *discordgo.Session, m *discordgo.MessageCreate) {
	intG, _ := strconv.Atoi(m.GuildID)
	intCh, _ := strconv.Atoi(m.ChannelID)
	if guild, ok := ext.Bot.Guilds[m.GuildID]; ok {
		//If guild already exists, just update the channel.
		guild.AdminChannel = intCh
		NoExcusesBot.WriteObject(ext.Bot.Database, guild)
		s.ChannelMessageSend(m.ChannelID, "This channel is now the admin channel!")
		ext.Bot.Logger.Printf("[CHANS] Updated admin channel in guild %v to %v", intG, intCh)
	} else {
		//If guild doesn't exist, create new guild
		guild = &NoExcusesBot.Guild{
			Id:           intG,
			AdminChannel: intCh,
		}
		ext.Bot.Guilds[m.GuildID] = guild
		NoExcusesBot.WriteObject(ext.Bot.Database, guild)
		s.ChannelMessageSend(m.ChannelID, "This channel is now the admin channel!")
		ext.Bot.Logger.Printf("[CHANS] Created admin channel in new guild %v to %v", intG, intCh)
	}
}
