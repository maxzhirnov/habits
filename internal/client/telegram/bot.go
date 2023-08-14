package telegram

import (
	"fmt"

	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"

	"github.com/maxzhirnov/habits/internal/client/telegram/handlers"
	"github.com/maxzhirnov/habits/internal/client/telegram/interfaces"
)

type Bot struct {
	APIHostName string
	*botapi.BotAPI
}

func NewBot(apiToken, apiHostName string) (*Bot, error) {
	bot, err := botapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}
	return &Bot{
		APIHostName: apiHostName,
		BotAPI:      bot,
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

	commands := map[string]interfaces.CommandHandler{
		"ping": &handlers.PingCommand{Bot: b},
		"add":  &handlers.AddNewHabit{Bot: b},
	}

	handler, exists := commands[update.Message.Command()]
	if !exists {
		log.Warnf("Command not supported: %s", update.Message.Command())
		return
	}

	msg, err := handler.Execute(update)
	if err != nil {
		log.Warn(err)
	}
	_, err = b.Send(msg)
	if err != nil {
		log.Warn(err)
	}
}

func (b *Bot) GetAPIURL(route string) string {
	return fmt.Sprintf("%s/%s", b.APIHostName, route)
}
