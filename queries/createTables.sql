PRAGMA foreign_keys = ON;
CREATE TABLE IF NOT EXISTS guilds (
    guildid INTEGER PRIMARY KEY, --Discord guild id
    adminrole INTEGER, --Discord role id for admins
    userchan INTEGER, --Discord channel id for user commands
    adminchan INTEGER --Discord channel id for admin commands
);
CREATE TABLE IF NOT EXISTS tasks (
    taskid INTEGER PRIMARY KEY AUTOINCREMENT,
    userid INTEGER NOT NULL, --Discord user id who created task
    guild INTEGER NOT NULL, --Refrences guilds.guildid for associated guild
    nextreminder INTEGER NOT NULL, --Unix timestamp of next reminder datetime
    interval INTEGER, --Interval of reminder in seconds (0 or NULL for disabled)
    repeats INTEGER, --How many times left to remind
    message TEXT, --Message of reminder
    FOREIGN KEY(guild) REFERENCES guilds(guildid)
);