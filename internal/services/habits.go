package services

import (
	"github.com/maxzhirnov/habits/internal/models"
	"github.com/maxzhirnov/habits/internal/utils"
	"time"
)

type Repository interface {
	CreateNewHabit(habit models.Habit) error
}

type App struct {
	Repo Repository
}

func NewAppService(repo Repository) *App {
	return &App{
		Repo: repo,
	}
}

func (app *App) AddNewHabit(h models.Habit) error {
	h.Tracking = map[time.Time]bool{
		utils.DateOnly(time.Now()): false,
	}
	if err := app.Repo.CreateNewHabit(h); err != nil {
		return err
	}
	return nil
}
