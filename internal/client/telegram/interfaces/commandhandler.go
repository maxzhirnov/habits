package interfaces

import botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type CommandHandler interface {
	Execute(update *botapi.Update) (*botapi.MessageConfig, error)
}
