package main

import (
	"database/sql"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const sqlUpdateGuildQuery = `
	UPDATE guilds 
	SET adminrole = $1, userchan = $2, adminchan = $3
	WHERE guildid = $4`
const sqlInsertGuildQuery = `
	INSERT INTO guilds (guildid, adminrole, userchan, adminchan)
	VALUES ($1, $2, $3, $4)`
const sqlInsertTaskQuery = `
	INSERT INTO tasks (userid, guild, nextreminder, interval, repeats, message)
	VALUES ($1, $2, $3, $4, $5, $6)`
const sqlDeleteTaskQuery = `
	DELETE FROM tasks
	WHERE taskid = $1`
const sqlUpdateTaskQuery = `
	UPDATE tasks
	SET nextreminder = $1, interval = $2, repeats = $3, message = $4
	WHERE taskid = $5`

type Guild struct {
	Id           int
	AdminRole    int
	UserChannel  int
	AdminChannel int
}

type Task struct {
	Id           int
	User         int
	Guild        int
	NextReminder int
	Interval     int
	Repeats      int
	Message      string
}

//Connects to Database and creates tables
func createDatabase() {
	database, _ = sql.Open("sqlite3", "tasks.db")
	dat, err := os.ReadFile("queries/createTables.sql")
	if err != nil {
		logger.Fatal("[SETUP] Error creating SQL tables")
	}
	database.Exec(string(dat))
}

//Get tasks from database
func getTasks() {
	tasks = make(map[int]Task)
	pendingTasks = make(map[int]Task)
	rows, err := database.Query("SELECT * FROM tasks")
	if err != nil {
		logger.Fatal("[SETUP] Error reading tasks from DB")
	}
	for rows.Next() {
		t := Task{}
		if err := rows.Scan(&t.Id, &t.User, &t.Guild, &t.NextReminder, &t.Interval, &t.Repeats, &t.Message); err != nil {
			logger.Fatalf("[SETUP] Could not scan task: %v", err)
		}
		tasks[t.Id] = t
	}
}

//Get guilds from database
func getGuilds() {
	guilds = make(map[int]Guild)
	rows, err := database.Query("SELECT * FROM guilds")
	if err != nil {
		logger.Fatal("[SETUP] Error reading tasks from DB")
	}
	for rows.Next() {
		g := Guild{}
		if err := rows.Scan(&g.Id, &g.AdminRole, &g.UserChannel, &g.AdminChannel); err != nil {
			logger.Fatalf("[SETUP] Could not scan guild: %v", err)
		}
		guilds[g.Id] = g
	}
}

//Updates guild to new info
func updateGuild(guildId int) {
	guild := guilds[guildId]
	_, err := database.Exec(sqlUpdateGuildQuery,
		guild.AdminRole,
		guild.UserChannel,
		guild.AdminChannel,
		guildId)
	if err != nil {
		logger.Fatalf("[CMD] Could not update guild: %v", err)
	}
}

//Create new guild
func createGuild(guildId int) {
	guild := guilds[guildId]
	_, err := database.Exec(sqlInsertGuildQuery,
		guildId,
		guild.AdminRole,
		guild.UserChannel,
		guild.AdminChannel)
	if err != nil {
		logger.Fatalf("[CMD] Could not create guild: %v", err)
	}
}

//Inserts a new task into the Database
func insertNewTask(t Task) int {
	res, err := database.Exec(sqlInsertTaskQuery,
		t.User,
		t.Guild,
		t.NextReminder,
		t.Interval,
		t.Repeats,
		t.Message)
	if err != nil {
		logger.Fatalf("[CMD] Could not create task: %v", err)
	}
	inserted, _ := res.LastInsertId()
	//Removing task from pending
	delete(pendingTasks, t.User)
	//Adding to internal task list
	t.Id = int(inserted)
	tasks[t.Id] = t
	return t.Id
}

//Deletes a task
func deleteTask(t Task) {
	//Deleting from DB
	_, err := database.Exec(sqlDeleteTaskQuery, t.Id)
	if err != nil {
		logger.Fatalf("[CMD] Could not delete task: %v", err)
	}
	//Deleting from memory
	delete(tasks, t.Id)
}

//Updates a task
func updateTask(t Task) {
	//Updating in DB
	_, err := database.Exec(sqlUpdateTaskQuery,
		t.NextReminder,
		t.Interval,
		t.Repeats,
		t.Message,
		t.Id)
	if err != nil {
		logger.Fatalf("[CMD] Could not update task: %v", err)
	}
	//Updating in memory
	tasks[t.Id] = t
}
