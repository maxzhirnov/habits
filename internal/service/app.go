package service

import (
	"github.com/maxzhirnov/habits/internal/models"
	"github.com/maxzhirnov/habits/internal/utils"
	"time"
)

type Repository interface {
	CreateNewHabit(habit models.Habit) error
	GetAllUsersHabits(userID string) ([]models.Habit, error)
	MarkHabitChecked(habit models.Habit) error
}

type App struct {
	Repo Repository
}

func NewApp(repo Repository) *App {
	return &App{
		Repo: repo,
	}
}

func (app App) AddNewHabit(h models.Habit) error {
	h.Streak = 0
	date := utils.DateOnly(time.Now())
	h.Tracking = map[time.Time]bool{
		date: false,
	}
	if err := app.Repo.CreateNewHabit(h); err != nil {
		return err
	}
	return nil
}

func (app App) GetAllUserHabits(userID string) ([]models.Habit, error) {
	return app.Repo.GetAllUsersHabits(userID)
}

func (app App) MarkHabitChecked(h models.Habit) error {
	return app.Repo.MarkHabitChecked(h)
}
