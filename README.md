# NoExcusesBot
This repository is for the community development of the NoExcuses Discord bot, it is written in the Go programming language.

## Installation
To run the bot locally, make sure you have the latest Go toolchain.

Clone the repository, download the requirements, and build the binary using:
```bash
git clone https://github.com/Potential-Central/NoExcusesBot.git
go get no-excuses
go build .
```
Make sure you have the environment variable exposed:
```
DISCORD_TOKEN=yourdiscordtoken
```
Or create a file called `.env` with the same contents.

## Running
Now just run the compiled binary according to your OS

`no-excuses.exe` in Windows, and `./no-excuses` in Linux and Mac

## User Commands
- **!time** - This command prints the current time in UTC
- **!newtask** - This command is responsible for creating a new task by parsing user arguments. The format is `!newtask <starttime> <interval> <repeats> <message>`. The arguments don't have to be supplied in this specific order, although the `<message>` should appear last.
    - `<starttime>`: This arguments specifies the starting time of the task and is provided in the format: 24:00 (If the starting time is earlier than the current time, the task will start tomorrow at this time)
    - `<interval>`: This arguments specifies the delay between tasks and is provided in the format 1w2d3h. This will translate into one week, two days, three hours. You can omit unneeded units: 2d is still a valid format.
    - `<repeats>`: This arguments specifies how many times should the task be repeated, this is just an integer. If you don't want the task to repeat, just provide 1.
    - `<message>`: This is just a message to be printed with the reminder, it's free form text.
- **!confirm** - This command is used to confirm that the task details are correct and it can start being executed.

## Admin Commands
- **!userchan** - This command sets the current channel to be the reminder channel - The channel where users can create tasks and that the reminders are sent in.
- **!adminchan** - This command sets the current channel to be the admin channel - The channel where moderators can use other administrative commands without users seeing.
- **!adminrole** - This command can only be used by a server admin, it recieves a role mention and sets that role to be able to execute the other admin commands.