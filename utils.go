package main

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	timeRegex     *regexp.Regexp
	intervalRegex *regexp.Regexp
	repeatsRegex  *regexp.Regexp
)

//Compiles regexes when starting the bot
func CompileRegex() {
	timeRegex, _ = regexp.Compile("\\d{1,2}:\\d{1,2}")
	intervalRegex, _ = regexp.Compile("(\\d{1,2}[dwh])+")
	repeatsRegex, _ = regexp.Compile("\\s\\d{1,3}\\s")
}

//Checks whether message author has administrator permissions
func IsAdmin(s *discordgo.Session, m *discordgo.MessageCreate) (bool, error) {
	//Get user permissions
	perm, err := s.UserChannelPermissions(m.Author.ID, m.ChannelID)
	if err != nil {
		return false, err
	}

	//Check for Guild admins
	if perm&discordgo.PermissionAdministrator != 0 {
		return true, nil
	}

	//Check for admin role in admin channel
	intCh, _ := strconv.Atoi(m.ChannelID)
	intG, _ := strconv.Atoi(m.GuildID)
	if guild, ok := guilds[intG]; ok {
		//Check for admin channel if one exists
		if intCh == guild.AdminChannel || guild.AdminChannel == 0 {
			//Check if user has admin role
			for _, role := range m.Member.Roles {
				intR, _ := strconv.Atoi(role)
				if intR == guild.AdminRole {
					return true, nil
				}
			}
		}
	}
	return false, nil
}

//Recieves arguments from a discord message
//Returns a Task according to given params
func ParseTaskArgs(arguments string) {
	//Parsing command
	msg := strings.Split(arguments, " ")[4:]
	logger.Println("MSG: ", msg)
	logger.Println("Time: ", timeRegex.FindString(arguments))
	logger.Println("Interval: ", intervalRegex.FindString(arguments))
	logger.Println("Repeats: ", repeatsRegex.FindString(arguments))
}
