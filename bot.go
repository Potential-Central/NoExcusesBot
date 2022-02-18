package NoExcusesBot

import (
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	bolt "go.etcd.io/bbolt"
)

// The basic bot object
type Bot struct {
	Prefix   string
	Token    string
	Logger   *log.Logger
	Session  *discordgo.Session
	Database *bolt.DB
	Guilds   map[string]*Guild
	Exts     []Extension
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
	ret := &Bot{Prefix: prefix, Token: token, Exts: make([]Extension, 0)}
	ret.Logger = log.Default()
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
		bot.Logger.Fatal("[SETUP] Error opening connection,", err)
	}
	bot.Logger.Println("[SETUP] Now running! Press CTRL+C to exit")
}

//Stops the bot.
func (bot *Bot) Stop() {
	bot.Logger.Println("[SETUP] Shutting down")
	bot.Session.Close()
	bot.Database.Close()
}

//Loads guilds from database, returns number of guilds
func (bot *Bot) LoadGuilds() int {
	bot.Guilds = make(map[string]*Guild)
	keys := GetKeys(bot.Database, "Guilds")
	for _, key := range keys {
		guildId, _ := strconv.Atoi(key)
		guild := &Guild{Id: guildId}
		ReadObject(bot.Database, guild)
		bot.Guilds[key] = guild
	}
	return len(keys)
}

func (bot *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Ignoring messages sent by self
	if m.Author.ID == s.State.User.ID {
		return
	}
	//Catching message by prefix
	if strings.HasPrefix(m.Content, bot.Prefix) {
		//Searching for commands in all extensions
		for _, ext := range bot.Exts {
			//Searching for all commands in each extension
			for _, cmd := range ext.CommandList() {
				if strings.HasPrefix(m.Content, bot.Prefix+cmd.Name) {
					//If command found, and user has permission, execute it
					if cmd.HasPermission(s, m) {
						bot.Logger.Printf("[COMND] Proccessing %s from ext %s.", cmd.Name, ext.Name())
						cmd.Execute(s, m)
						return
					}
				}
			}
		}
	}
}
