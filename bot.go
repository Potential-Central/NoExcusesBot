package NoExcusesBot

import (
	"log"
	"strconv"

	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

// The basic bot object
type Bot struct {
	Prefix   string
	Token    string
	logger   *log.Logger
	Session  *discordgo.Session
	Database *bolt.DB
	Guilds   map[string]Guild
}

// Guild object, implements DataObject
type Guild struct {
	Id           int `json:"id"`
	AdminRole    int `json:"adminRole,omitempty"`
	UserChannel  int `json:"userChannel,omitempty"`
	AdminChannel int `json:"adminChannel,omitempty"`
}

func (guild Guild) bucket() string {
	return "Guilds"
}

func (guild Guild) primaryKey() string {
	return strconv.Itoa(guild.Id)
}

// Takes a token and creates a discord session, initializes the bot.
// 2nd argument is an optional prefix, default is !
func MakeBot(token string, args ...string) (*Bot, error) {
	var err error
	prefix := "!"
	if len(args) >= 1 {
		prefix = args[0]
	}
	ret := &Bot{Prefix: prefix, Token: token}
	ret.logger = log.Default()
	ret.Session, err = discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	ret.Database, err = CreateDB("Guilds")
	if err != nil {
		return nil, err
	}
	ret.LoadGuilds()
	ret.Session.Identify.Intents = discordgo.IntentsGuilds | discordgo.IntentsGuildMembers | discordgo.IntentsGuildMessages
	ret.Session.AddHandler(ret.messageHandler)
	return ret, nil
}

//Starts the bot, this function is not blocking.
func (bot *Bot) Start() {
	err := bot.Session.Open()
	if err != nil {
		bot.logger.Fatal("[SETUP] Error opening connection,", err)
	}
	bot.logger.Println("[SETUP] Now running! Press CTRL+C to exit")
}

//Stops the bot.
func (bot *Bot) Stop() {
	bot.logger.Println("[SETUP] Shutting down")
	bot.Session.Close()
	bot.Database.Close()
}

//Loads guilds from database, returns number of guilds
func (bot *Bot) LoadGuilds() int {
	bot.Guilds = make(map[string]Guild)
	keys := GetKeys(bot.Database, "Guilds")
	for _, key := range keys {
		guildId, _ := strconv.Atoi(key)
		guild := &Guild{Id: guildId}
		ReadObject(bot.Database, guild)
		bot.Guilds[key] = *guild
	}
	return len(keys)
}

func (bot *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	bot.logger.Println("Test")
}
