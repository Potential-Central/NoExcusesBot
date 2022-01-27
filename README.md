# NoExcusesBot
This repository is for the community development of the NoExcuses Discord bot, it is written in the Go programming language.

## Main Idea
The bot uses UTC timezone, allows registration of personalized reminders by memebers of the group. After a reminder is sent, the bot will ask for confirmation of task completion.

## Feature List
- Task creation, and deletion system with intervals.
- Reminder scheduling.
- Confirmation system.
- Admin control for reminders.

## Thoughts
- Task creation should probably be interactive rather than a one message thing, ideally will default to interactive mode if no arguments specified.
- Tasks should be saved in SQL, (ID, User, NextTimestamp, Interval, Repeats, Message).
- Every set interval (2 minutes or something) bot will query all tasks with NextTimestamp < Now. Then notify in chat, and set NextTimestamp+=Interval.
- After reminder is sent, decrement the Reapets by 1.
- Checking for confirmation might be tricky, need to save in memory message IDs of sent reminders and associated users. 
- Go's Built-in time parser doesn't support days and weeks, consider using [go-str2duration](https://github.com/xhit/go-str2duration).
- Minimum interval should be established to avoid abuse.