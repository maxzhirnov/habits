package handlers

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhirnov/habits/internal/client/telegram/interfaces"
)

type PingCommand struct {
	Bot interfaces.TelegramBot
}

func (c *PingCommand) Execute(update *botapi.Update) (*botapi.MessageConfig, error) {
	msg := botapi.NewMessage(update.Message.From.ID, "pong")
	return &msg, nil
}
