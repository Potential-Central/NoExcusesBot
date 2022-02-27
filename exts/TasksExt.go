package NoExcusesBot

import (
	"strconv"

	"github.com/Potential-Central/NoExcusesBot"
	"github.com/bwmarrin/discordgo"
)

// Task object implementing DataObject
type Task struct {
	Id           int    `json:"id"`
	Guild        int    `json:"guild,omitempty"`
	User         int    `json:"user,omitempty"`
	NextReminder int    `json:"nextReminder,omitempty"`
	Interval     int    `json:"interval,omitempty"`
	Repeats      int    `json:"repeats,omitempty"`
	Message      string `json:"message,omitempty"`
}

func (task *Task) Bucket() string {
	return "Tasks"
}

func (task *Task) PrimaryKey() string {
	return strconv.Itoa(task.Id)
}

type TasksExt struct {
	Bot      *NoExcusesBot.Bot
	Commands []*NoExcusesBot.Command
	Tasks    map[int]*Task
}

func MakeTasksExt(bot *NoExcusesBot.Bot) {
	ret := &TasksExt{bot, make([]*NoExcusesBot.Command, 0), make(map[int]*Task)}

	//Loading tasks
	ret.loadTasks()

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

func (ext *TasksExt) loadTasks() {
	for _, key := range NoExcusesBot.GetKeys(ext.Bot.Database, "Tasks") {
		id, _ := strconv.Atoi(key)
		task := &Task{Id: id}
		NoExcusesBot.ReadObject(ext.Bot.Database, task)
		ext.Tasks[id] = task
	}
}
