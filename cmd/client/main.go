package main

import (
	"fmt"
	"github.com/maxzhirnov/habits/internal/client/telegram"
	"github.com/maxzhirnov/habits/internal/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.New()
	cfg.Parse()

	apiUrl := fmt.Sprintf("http://%s:%d", cfg.Client.Host, cfg.Client.Port)
	bot, err := telegram.NewBot(cfg.Client.TelegramToken, apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	bot.Run()
}
