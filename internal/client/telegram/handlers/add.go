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
		var (
			res       = ""
			url       = apiUrl + "add-new-habit"
			habitName = c.Message().Payload
			userID    = c.Sender().ID
			date      = c.Message().Time()
			client    = &http.Client{}
			data      = []byte(fmt.Sprintf(`{"name": "%s", "user_id": "%d", "created": "%s"}`, habitName, userID, date))
		)

		//If the command payload is empty giving user help message
		if habitName == "" {
			res = "To add new habit use \"/add habit name\", for example \"/add running every day\""
			return c.Send(res)
		}

		//Trying to sava new habit item via API call
		req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		//Handling response codes
		switch resp.StatusCode {
		default:
			//TODO: should something be sent to the user here?
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
	}
}
