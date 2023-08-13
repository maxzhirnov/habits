package repos

import (
	"context"
	"errors"
	"github.com/maxzhirnov/habits/internal/models"
	"time"
)

var ErrHabitExists = errors.New("habit already exists")

type DataBase interface {
	Insert(ctx context.Context, collectionName string, document interface{}) (string, error)
	Exists(ctx context.Context, collectionName string, filters map[string]interface{}) (bool, error)
}

type Repository struct {
	DB DataBase
}

func New(db DataBase) *Repository {
	return &Repository{
		DB: db,
	}
}

func (r Repository) CreateNewHabit(habit models.Habit) error {
	collectionName := "habits"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := map[string]interface{}{
		"name":    habit.Name,
		"user_id": habit.UserID,
	}
	exists, err := r.DB.Exists(ctx, collectionName, filter)
	if err != nil {
		return err
	}

	if exists {
		return ErrHabitExists
	}

	_, err = r.DB.Insert(ctx, collectionName, habit)
	if err != nil {
		return err
	}
	return nil
}
