package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/maxzhirnov/habits/internal/models"
	"net/http"
)

func getHabits(userID int64, apiURL string) ([]models.Habit, error) {
	var (
		url    = fmt.Sprintf("%s%s?userid=%d", apiURL, "habits", userID)
		client = &http.Client{}
	)

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch habits, status code: %d", resp.StatusCode)
	}

	var habits []models.Habit
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&habits)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}
	return habits, err
}
