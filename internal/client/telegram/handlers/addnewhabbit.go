package handlers

import (
	"bytes"
	"fmt"
	botapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

type AddNewHabit struct {
	APIURL string
}

func (c *AddNewHabit) Execute(update *botapi.Update) (string, error) {
	url := c.APIURL + "add-new-habit"
	res := ""

	habitName := update.Message.Text
	habitName = strings.TrimPrefix(habitName, "/add")
	habitName = strings.TrimSpace(habitName)

	if habitName == "" {
		res = "To add new habit use \"/add habit name\", for example \"/add running every day\""
		return res, nil
	}

	userID := strconv.FormatInt(update.Message.From.ID, 10)

	// Создание JSON-тела запроса
	data := []byte(fmt.Sprintf(`{
		"name": "%s",
		"user_id": "%s"
	}`, habitName, userID))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}

	// Установка заголовка Content-Type
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	default:
		return "", fmt.Errorf("unexpected response status: %s", resp.Status)
	case http.StatusOK:
		res = fmt.Sprintf("New habit added: \"%s\"", habitName)
		log.Info(res)
	case http.StatusConflict:
		res = fmt.Sprintf("Habit \"%s\" already exists", habitName)
		log.Warn(res)
	case http.StatusBadRequest:
		res = fmt.Sprint("Something went wrong")
		log.Warn("Bad Request: 400")
	}

	return res, nil
}
