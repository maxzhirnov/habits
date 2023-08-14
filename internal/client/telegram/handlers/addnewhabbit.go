package handlers

import (
	"bytes"
	"fmt"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/maxzhirnov/habits/internal/client/telegram/interfaces"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type AddNewHabit struct {
	Bot interfaces.TelegramBot
}

func (c *AddNewHabit) Execute(update *botapi.Update) (*botapi.MessageConfig, error) {
	url := c.Bot.GetAPIURL("add-new-habit")

	habitName := update.Message.Text
	habitName = strings.TrimPrefix(habitName, "/add")
	habitName = strings.TrimSpace(habitName)
	userID := strconv.FormatInt(update.Message.From.ID, 10)

	// Создание JSON-тела запроса
	data := []byte(fmt.Sprintf(`{
		"name": "%s",
		"user_id": "%s"
	}`, habitName, userID))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	msg := botapi.NewMessage(update.Message.From.ID, "")
	switch resp.StatusCode {
	default:
		return nil, fmt.Errorf("unexpected response status: %s", resp.Status)
	case http.StatusOK:
		msg.Text = fmt.Sprintf("New habit added: \"%s\"", habitName)
		log.Info(msg.Text)
	case http.StatusConflict:
		msg.Text = fmt.Sprintf("Habit \"%s\" already exists", habitName)
		log.Warn(msg.Text)
	case http.StatusBadRequest:
		msg.Text = fmt.Sprint("Something went wrong")
		log.Warn("Bad Request: 400")
	}

	return &msg, nil
}
