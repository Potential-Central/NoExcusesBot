package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

//Cron job for checking tasks
func CheckTasks() {
	now := time.Now().UTC()
	logger.Printf("[TASK] Checking reminders")
	for i := range tasks {
		//Checking if reminder is overdue
		reminder_time := time.Unix(int64(tasks[i].NextReminder), 0).UTC()
		if reminder_time.Before(now) {
			HandleTask(i)
		}
	}
}

//Handles a single task
func HandleTask(t int) {
	task, _ := tasks[t]
	task = SendTask(task)
	//If no more repeats, delete task
	if task.Repeats <= 0 {
		logger.Printf("[TASK] Deleting Task %v", task.Id)
		deleteTask(task)
	} else {
		logger.Printf("[TASK] Updating Task %v", task.Id)
		updateTask(task)
	}
}

func SendTask(t Task) Task {
	//Rescheduling task
	t.Repeats -= 1
	t.NextReminder += t.Interval
	//Mentioning creator
	userMention := fmt.Sprintf("<@%v> %s", t.User, t.Message)
	embed := TaskToEmbed(t, "Task Reminder!", "", 16112962)
	//Sending
	if g, ok := guilds[t.Guild]; ok {
		_, err := client.ChannelMessageSendComplex(strconv.Itoa(g.UserChannel), &discordgo.MessageSend{Content: userMention, Embed: &embed, TTS: false})
		if err != nil {
			logger.Printf("[TASK] Error sending reminder! %v", err)
		} else {
			logger.Printf("[TASK] Sending reminder! %v", t.Id)
		}
	} else {
		logger.Println("[TASK] Error sending reminder! Guild not found")
	}
	return t
}
