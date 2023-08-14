package telegram

import (
	"github.com/maxzhirnov/habits/internal/client/telegram/handlers"
	log "github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"time"
)

type TeleBot struct {
	ApiUrl string
	*tele.Bot
}

func NewTeleBot(token, apiUrl string) (*TeleBot, error) {
	pref := tele.Settings{
		Token:  token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	return &TeleBot{
		ApiUrl: apiUrl,
		Bot:    b,
	}, nil
}

func (b TeleBot) Run() {
	b.Handle("/ping", handlers.Hello)
	b.Handle("/add", handlers.AddNewHabit(b.ApiUrl))

	b.Start()
}
