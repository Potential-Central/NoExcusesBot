package main

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

//Handles messages
func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Ignoring message bot sent
	if m.Author.ID == s.State.User.ID {
		return
	}

	//Admin commands
	if admin, _ := IsAdmin(s, m); admin {
		if strings.HasPrefix(m.Content, "!userchan") {
			cmdUserChan(s, m)
			return
		} else if strings.HasPrefix(m.Content, "!adminchan") {
			cmdAdminChan(s, m)
			return
		} else if strings.HasPrefix(m.Content, "!adminrole ") {
			cmdAdminRole(s, m)
			return
		}
	}
}

//Handles guild joins
func guildJoinHandler(s *discordgo.Session, g *discordgo.GuildCreate) {

}

//Setting a channel where reminder get posted, and users interact with the bot
func cmdUserChan(s *discordgo.Session, m *discordgo.MessageCreate) {
	intG, _ := strconv.Atoi(m.GuildID)
	intCh, _ := strconv.Atoi(m.ChannelID)
	if guild, ok := guilds[intG]; ok {
		//If guild already exists, just update the channel.
		guild.UserChannel = intCh
		guilds[intG] = guild
		updateGuild(intG)
		logger.Printf("[CMD] Updated user channel in guild %v", intG)
	} else {
		//If guild doesn't exist, create new guild
		newGuild := Guild{
			Id:          intG,
			UserChannel: intCh,
		}
		guilds[intG] = newGuild
		createGuild(intG)
		logger.Printf("[CMD] Created user channel in guild %v", intG)
	}
}

//Setting a channel where administrative commands are used
func cmdAdminChan(s *discordgo.Session, m *discordgo.MessageCreate) {
	intG, _ := strconv.Atoi(m.GuildID)
	intCh, _ := strconv.Atoi(m.ChannelID)
	if guild, ok := guilds[intG]; ok {
		//If guild already exists, just update the channel.
		guild.AdminChannel = intCh
		guilds[intG] = guild
		updateGuild(intG)
		logger.Printf("[CMD] Updated admin channel in guild %v", intG)
	} else {
		//If guild doesn't exist, create new guild
		newGuild := Guild{
			Id:           intG,
			AdminChannel: intCh,
		}
		guilds[intG] = newGuild
		createGuild(intG)
		logger.Printf("[CMD] Created admin channel in guild %v", intG)
	}
}

//Setting a role for bot administrators
func cmdAdminRole(s *discordgo.Session, m *discordgo.MessageCreate) {
	roles := m.MentionRoles
	intG, _ := strconv.Atoi(m.GuildID)
	if len(roles) != 1 {
		s.ChannelMessageSend(m.ChannelID, "Please provide exactly **one** role!")
		return
	}
	intR, _ := strconv.Atoi(roles[0])
	if guild, ok := guilds[intG]; ok {
		//If guild already exists, just update the role.
		guild.AdminRole = intR
		guilds[intG] = guild
		updateGuild(intG)
		logger.Printf("[CMD] Updated admin role in guild %v", intG)
		s.ChannelMessageSend(m.ChannelID, "Admin role successfully updated!")
	} else {
		//If guild doesn't exist, create new guild
		newGuild := Guild{
			Id:        intG,
			AdminRole: intR,
		}
		guilds[intG] = newGuild
		createGuild(intG)
		logger.Printf("[CMD] Created admin role in guild %v", intG)
		s.ChannelMessageSend(m.ChannelID, "Admin role successfully created!")
	}
}