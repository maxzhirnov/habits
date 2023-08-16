package handlers

import (
	"fmt"
	"github.com/maxzhirnov/habits/internal/models"
	"github.com/maxzhirnov/habits/internal/utils"
	log "github.com/sirupsen/logrus"
	tele "gopkg.in/telebot.v3"
	"strings"
	"time"
)

func ListHabits(apiUrl string) func(c tele.Context) error {
	return func(c tele.Context) error {
		habits, err := getHabits(c.Sender().ID, apiUrl)
		if err != nil {
			log.Warn(err)
			return err
		}
		if len(habits) == 0 {
			return c.Send("You don't have any habits to track yet. You can add one using \"/add habbit\" command")
		}

		return c.Send(habitsToString(habits))
	}
}

func habitsToString(habits []models.Habit) string {
	var habitStrings []string

	for _, h := range habits {

		today := utils.DateOnly(time.Now())
		totalCount := utils.DaysBetween(h.CratedAt, today) + 1
		trueCount := len(h.Tracking)

		lastDays := getLastDaysCheckboxes(h, 14)

		habitString := fmt.Sprintf(
			"%s\n%d/%d %s",
			h.Name,
			trueCount,
			totalCount,
			lastDays,
		)
		habitStrings = append(habitStrings, habitString)
	}

	return strings.Join(habitStrings, "\n\n")
}

func getLastDaysCheckboxes(h models.Habit, days int) string {
	tracking := h.Tracking
	var checkboxes []string
	now := time.Now()

	for i := 0; i < days; i++ {
		day := now.AddDate(0, 0, -i)
		dayStart := utils.DateOnly(day)

		dayCreated := utils.DateOnly(h.CratedAt)
		// Проверяем значение для этой даты
		checked, exists := tracking[dayStart]
		switch {
		case exists && checked:
			checkboxes = append([]string{"☑"}, checkboxes...) // prepend
		case dayCreated.After(dayStart) || dayCreated.Equal(dayStart):
			checkboxes = append([]string{"☐"}, checkboxes...)
		case i == 0 && !checked:
			checkboxes = append([]string{"☒"}, checkboxes...)
		case exists && !checked:
			checkboxes = append([]string{"☒"}, checkboxes...) // prepend
		default:
			checkboxes = append([]string{"☐"}, checkboxes...) // prepend
		}
	}

	return strings.Join(checkboxes, "")
}
