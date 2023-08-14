package telegram

import (
	"fmt"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/maxzhirnov/habits/internal/client/telegram/handlers"
)

type CommandHandler interface {
	Execute(update *botapi.Update) (string, error)
}

type Bot struct {
	APIURL   string
	Commands map[string]CommandHandler
	*botapi.BotAPI
}

func NewBot(tgToken, APIURL string) (*Bot, error) {
	bot, err := botapi.NewBotAPI(tgToken)
	if err != nil {
		return nil, err
	}

	commands := map[string]CommandHandler{
		"ping": &handlers.PingCommand{},
		"add":  &handlers.AddNewHabit{APIURL: APIURL},
	}

	return &Bot{
		Commands: commands,
		BotAPI:   bot,
	}, nil
}

func (b *Bot) Run() {
	updateConfig := botapi.NewUpdate(0)
	updateConfig.Timeout = 30

	updates := b.GetUpdatesChan(updateConfig)
	for update := range updates {
		b.handleUpdate(&update)
	}
}

func (b *Bot) handleUpdate(update *botapi.Update) {
	if update.Message == nil || !update.Message.IsCommand() {
		return
	}

	msg := botapi.NewMessage(update.Message.From.ID, "")
	var err error

	handler, exists := b.Commands[update.Message.Command()]
	if !exists {
		msg.Text = fmt.Sprintf("Command %s not supported, please use: %s instead", update.Message.Text, getCommandsString(b.Commands))
		log.Warnf("Command %s not supported: %s", update.Message.Text, update.Message.Command())
	} else {
		msg.Text, err = handler.Execute(update)
		if err != nil {
			log.Warn(err)
			return
		}
	}

	_, err = b.Send(msg)
	if err != nil {
		log.Warn(err)
		return
	}
}
