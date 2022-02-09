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
			} else if strings.HasPrefix(m.Content, "!confirm") {
				cmdTaskConfirm(s, m)
				return
			} else if strings.HasPrefix(m.Content, "!clock") {
				cmdClock(s, m)
				return
			} else if strings.HasPrefix(m.Content, "!time") {
				cmdClock(s, m)
				return
			} else if strings.HasPrefix(m.Content, "!help") {
				cmdHelp(s, m)
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
	//Making sure we have enough arguments
	if len(strings.Split(m.Content, " ")) < 5 {
		logger.Printf("[CMD] Invalid task creation %v", m.Content)
		s.ChannelMessageSend(m.ChannelID, "`Invalid task creation! Use the format: !newtask starttime interval repeats message`")
		return
	}
	//Parsing
	task, err := ParseTaskArgs(m.Content)
	if err != nil {
		logger.Printf("[CMD] Invalid task creation %v", m.Content)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`%v`", err))
		return
	}
	task.User, _ = strconv.Atoi(m.Author.ID)
	task.Guild, _ = strconv.Atoi(m.GuildID)
	embed := TaskToEmbed(task, "Task Creation", "Type **!confirm** to create this task", 16106050)
	_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Embed: &embed, TTS: false})
	if err != nil {
		logger.Printf("[CMD] Error sending task: %v", err)
		return
	}
	//Appending task to pending, waiting for it to be confirmed via !confirm
	pendingTasks[task.User] = task
}

//Confirms a pending task, submits it to database
func cmdTaskConfirm(s *discordgo.Session, m *discordgo.MessageCreate) {
	userid, _ := strconv.Atoi(m.Author.ID)
	//Checking if user has a pending task
	if task, ok := pendingTasks[userid]; ok {
		task.Id = insertNewTask(task)
		embed := TaskToEmbed(task, fmt.Sprintf("Task **%v** successfully created!", task.Id), "", 261131)
		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Embed: &embed, TTS: false})
		if err != nil {
			logger.Printf("[CMD] Error sending task: %v", err)
			return
		}
	}
	logger.Printf("[CMD] Created task: %v", m.Content)
}

//Returns current time in UTC
func cmdClock(s *discordgo.Session, m *discordgo.MessageCreate) {
	cur := time.Now().UTC()
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("The current time is **%s**", cur.Format("15:04")))
}

//Prints help about given command
func cmdHelp(s *discordgo.Session, m *discordgo.MessageCreate) {
	if seperated := strings.Split(m.Content, " "); len(seperated) == 1 {
		//Main help page
		embed := help["help"]
		embed.Fields = make([]*discordgo.MessageEmbedField, 0, len(help))
		for k, v := range help {
			if k != "help" {
				embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
					Name:   fmt.Sprintf("!%s", k),
					Value:  v.Description,
					Inline: false,
				})
			}
		}
		_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Embed: &embed, TTS: false})
		if err != nil {
			logger.Printf("[CMD] Error sending help: %v", err)
		}
	} else if len(seperated) >= 2 {
		//Help for specific command
		keyword := strings.ToLower(strings.Replace(seperated[1], "!", "", 1))
		if embed, ok := help[keyword]; ok {
			_, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{Embed: &embed, TTS: false})
			if err != nil {
				logger.Printf("[CMD] Error sending help: %v", err)
			}
		} else {
			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Did not find the subject `%s` :( Try `!help` to see all commands available.", keyword))
			if err != nil {
				logger.Printf("[CMD] Error sending help: %v", err)
			}
		}
	}
}
