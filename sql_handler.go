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
	guild, _ := guilds[guildId]
	_, err := database.Exec(sqlUpdateGuildQuery,
		guild.AdminRole,
		guild.UserChannel,
		guild.AdminChannel,
		guildId)
	if err != nil {
		logger.Fatalf("[CMD] Could not update guild: %v", err)
	}
	logger.Printf("[CND] Updated guild %v", guildId)
}

//Create new guild
func createGuild(guildId int) {
	guild, _ := guilds[guildId]
	_, err := database.Exec(sqlInsertGuildQuery,
		guildId,
		guild.AdminRole,
		guild.UserChannel,
		guild.AdminChannel)
	if err != nil {
		logger.Fatalf("[CMD] Could not create guild: %v", err)
	}
	logger.Printf("[CND] Created guild %v", guildId)
}
