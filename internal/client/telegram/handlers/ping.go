package handlers

import (
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type PingCommand struct {
}

func (c *PingCommand) Execute(update *botapi.Update) (string, error) {
	return "pong", nil
}
