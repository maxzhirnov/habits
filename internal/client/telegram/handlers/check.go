package handlers

import (
	"bytes"
	"fmt"
	"github.com/maxzhirnov/habits/internal/models"
	tele "gopkg.in/telebot.v3"
	"net/http"
	"strings"
)

func TrackHabitButtons(apiUrl string) func(c tele.Context) error {
	return func(c tele.Context) error {
		m := c.Message()
		habits, err := getHabits(c.Sender().ID, apiUrl)
		if err != nil {
			return err
		}
		if len(habits) == 0 {
			return c.Send("You don't have any habits to track yet. You can add one using \"/add habbit\" command")
		}

		buttons, hasUndone := generateButtons(habits)

		if hasUndone == false {
			return c.Send("All habits for today already marked as done. Good job!")
		}

		_, err = c.Bot().Send(m.Sender, "Choose habit you want to mark as checked for today:", buttons)
		if err != nil {
			return err
		}
		return nil
	}
}

func generateButtons(data []models.Habit) (*tele.ReplyMarkup, bool) {
	var buttons [][]tele.InlineButton
	for _, item := range data {
		if !item.IsDoneToday() {
			btn := tele.InlineButton{
				Unique: "btn_" + item.Name,
				Text:   item.Name,
				Data:   item.ID,
			}
			buttons = append(buttons, []tele.InlineButton{btn})
		}
	}
	if len(buttons) == 0 {
		return nil, false
	}

	markup := &tele.ReplyMarkup{InlineKeyboard: buttons}
	return markup, true
}

func TrackHabitCallback(apiUrl string) func(tele.Context) error {
	return func(c tele.Context) error {
		callbackDataCleaned := strings.TrimSpace(c.Callback().Data)

		if strings.HasPrefix(callbackDataCleaned, "btn_") {
			var (
				objID  = strings.Split(callbackDataCleaned, "|")[1]
				url    = fmt.Sprintf("%s%s", apiUrl, "mark")
				client = &http.Client{}
				data   = []byte(fmt.Sprintf(`{"id": "%s", "user_id": "%d"}`, objID, c.Sender().ID))
			)

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

			if resp.StatusCode != http.StatusOK {
				return c.Send("Something went wrong")
			}
			err = c.Delete()
			if err != nil {
				return err
			}
			return ListHabits(apiUrl)(c)
		}
		return nil
	}
}
