package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
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
		} else if strings.HasPrefix(m.Content, "!adminrole") {
			cmdAdminRole(s, m)
			return
		}
	}
	//User commands
	intG, _ := strconv.Atoi(m.GuildID)
	if guild, ok := guilds[intG]; ok {
		intC, _ := strconv.Atoi(m.ChannelID)
		if intC == guild.UserChannel || intC == guild.AdminChannel {
			if strings.HasPrefix(m.Content, "!newtask") {
				cmdNewTask(s, m)
				return
			} else if strings.HasPrefix(m.Content, "!clock") {
				cmdClock(s, m)
				return
			} else if strings.HasPrefix(m.Content, "!time") {
				cmdClock(s, m)
				return
			}
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
		s.ChannelMessageSend(m.ChannelID, "This channel is now the reminder channel!")
	} else {
		//If guild doesn't exist, create new guild
		newGuild := Guild{
			Id:          intG,
			UserChannel: intCh,
		}
		guilds[intG] = newGuild
		createGuild(intG)
		logger.Printf("[CMD] Created user channel in guild %v", intG)
		s.ChannelMessageSend(m.ChannelID, "This channel is now the reminder channel!")
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
		s.ChannelMessageSend(m.ChannelID, "This channel is now an admin channel!")
	} else {
		//If guild doesn't exist, create new guild
		newGuild := Guild{
			Id:           intG,
			AdminChannel: intCh,
		}
		guilds[intG] = newGuild
		createGuild(intG)
		logger.Printf("[CMD] Created admin channel in guild %v", intG)
		s.ChannelMessageSend(m.ChannelID, "This channel is now an admin channel!")
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

//Creating a new task
//Command format: !newtask starttime interval repeats message
//Command example: !newtask 15:00 1d 1 Did you remember to plant your radishes?
func cmdNewTask(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(strings.Split(m.Content, " ")) < 5 {
		logger.Printf("[CMD] Invalid task creation %v", m.Content)
		s.ChannelMessageSend(m.ChannelID, "`Invalid task creation! Use the format: !newtask starttime interval repeats message`")
		return
	}
	ParseTaskArgs(m.Content)
}

//Returns current time in UTC
func cmdClock(s *discordgo.Session, m *discordgo.MessageCreate) {
	cur := time.Now().UTC()
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The current time is **%s**", cur.Format("15:04")))
}
