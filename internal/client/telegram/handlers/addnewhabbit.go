package handlers

import (
	"bytes"
	"fmt"
	log "github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"net/http"
)

func AddNewHabit(apiUrl string) func(c tele.Context) error {
	return func(c tele.Context) error {
		url := apiUrl + "add-new-habit"
		res := ""

		habitName := c.Message().Payload
		//habitName = strings.TrimPrefix(habitName, "/add")
		//habitName = strings.TrimSpace(habitName)

		if habitName == "" {
			res = "To add new habit use \"/add habit name\", for example \"/add running every day\""
			return c.Send(res)
		}

		userID := c.Sender().ID

		// Создание JSON-тела запроса
		data := []byte(fmt.Sprintf(`{
		"name": "%s",
		"user_id": "%d"
	}`, habitName, userID))

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			return err
		}

		// Установка заголовка Content-Type
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		switch resp.StatusCode {
		default:
			return fmt.Errorf("unexpected response status: %s", resp.Status)
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

		return c.Send(res)

		return c.Send("Pong")
	}
}
