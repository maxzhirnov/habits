package main

import (
	"fmt"
	"github.com/maxzhirnov/habits/internal/client/telegram"
	"github.com/maxzhirnov/habits/internal/config"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	time.Local = time.UTC
	log.Infof("Starting. Local time: %s", time.Now())
	cfg := config.New()
	cfg.Parse()

	apiUrl := fmt.Sprintf("http://%s:%d/", cfg.Client.Host, cfg.Client.Port)
	bot, err := telegram.NewTeleBot(cfg.Client.TelegramToken, apiUrl)
	if err != nil {
		log.Fatal(err)
	}

	bot.Run()
}
