package NoExcusesBot

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// The basic bot object
type Bot struct {
	Prefix  string
	Token   string
	logger  *log.Logger
	Session *discordgo.Session
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
}

//Stops the bot.
func (bot *Bot) Stop() {
	bot.Session.Close()
}

func (bot *Bot) messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	bot.logger.Println("Test")
}
