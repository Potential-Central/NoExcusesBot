package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	str2duration "github.com/xhit/go-str2duration/v2"
)

var (
	timeRegex     *regexp.Regexp
	intervalRegex *regexp.Regexp
	repeatsRegex  *regexp.Regexp
)

//Seriously, get help
func GetHelp() {
	help = make(map[string]discordgo.MessageEmbed)
	jsonFile, err := os.Open("help.json")
	if err != nil {
		return
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &help)
}

//Compiles regexes when starting the bot
func CompileRegex() {
	timeRegex, _ = regexp.Compile("\\d{1,2}:\\d{1,2}")
	intervalRegex, _ = regexp.Compile("([\\d\\.]{1,3}[dwh])+")
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
func ParseTaskArgs(arguments string) (Task, error) {
	//Parsing command
	cur := time.Now().UTC()
	msgArg := strings.Join(strings.Split(arguments, " ")[4:], " ")
	//Getting time in the HH:MM format
	timeArg, err := time.Parse("15:04", timeRegex.FindString(arguments))
	if err != nil {
		return Task{}, errors.New("Time format not recognized; Use 24:00 format.")
	}
	//Geting interval of reminders in the 1w2d3h format
	intervalArg, err := str2duration.ParseDuration(intervalRegex.FindString(arguments))
	if err != nil || intervalArg <= time.Hour*1 {
		return Task{}, errors.New("Error with interval; Please only use (h)ours (d)ays and (w)eeks. Must be more than 1 hour.")
	}
	//Getting number of repeats
	repeatArg, err := strconv.Atoi(strings.Trim(repeatsRegex.FindString(arguments), " "))
	if err != nil {
		return Task{}, errors.New("Error with repeats; Please specify number of times to repeat reminder. If you don't want to repeat, Use 1")
	}
	//Next time should be the time given in the same date as today,
	//Unless time is already in the past, then time should be tomorrow.
	next := time.Date(cur.Year(), cur.Month(), cur.Day(), timeArg.Hour(), timeArg.Minute(), 0, 0, time.UTC)
	if next.Before(cur) {
		next = next.Add(time.Hour * 24)
	}
	return Task{
		NextReminder: int(next.Unix()),
		Interval:     int(intervalArg.Seconds()),
		Repeats:      repeatArg,
		Message:      msgArg}, nil
}

//Recieves a task and constructs an embed for it
//TODO: make this less ugly...
func TaskToEmbed(t Task, title, desc string, color int) discordgo.MessageEmbed {
	nxt := time.Unix(int64(t.NextReminder), 0).UTC().Format("02-01-2006 15:04")
	nxtDiff := int(time.Unix(int64(t.NextReminder), 0).UTC().Sub(time.Now().UTC()).Hours())
	embed := discordgo.MessageEmbed{
		Title:       title,
		Description: desc,
		Color:       color,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "Next Reminder",
				Value: fmt.Sprintf("%s (in %v hours)", nxt, nxtDiff),
			},
			{
				Name:  "Interval",
				Value: (time.Duration(t.Interval) * time.Second).String(),
			},
			{
				Name:  "Repeats",
				Value: fmt.Sprintf("Repeats %v times", t.Repeats),
			},
			{
				Name:  "Message",
				Value: t.Message,
			},
		},
	}
	return embed
}
