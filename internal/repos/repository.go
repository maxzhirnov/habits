package repos

import (
	"context"
	"errors"
	"fmt"
	"github.com/maxzhirnov/habits/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

var (
	ErrHabitExists   = errors.New("habit already exists")
	ErrHabitNotFound = errors.New("habit not found")
)

type DataBase interface {
	Insert(ctx context.Context, collectionName string, document interface{}) (string, error)
	Exists(ctx context.Context, collectionName string, filters map[string]interface{}) (bool, error)
	GetAll(ctx context.Context, collectionName string, filter map[string]interface{}) ([]interface{}, error)
	UpdateHabitTracking(ctx context.Context, collectionName string, habitID string, date time.Time, value bool) error
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

func (r Repository) GetAllUsersHabits(userID string) ([]models.Habit, error) {
	filter := map[string]interface{}{"user_id": userID}
	items, err := r.DB.GetAll(context.Background(), "habits", filter)
	if err != nil {
		return nil, err
	}

	habits := make([]models.Habit, 0)
	for _, item := range items {
		rawBson, ok := item.(bson.D)
		if !ok {
			return nil, fmt.Errorf("failed to convert item to bson.M")
		}

		bytes, err := bson.Marshal(rawBson)
		if err != nil {
			return nil, err
		}

		var habit models.Habit
		err = bson.Unmarshal(bytes, &habit)
		if err != nil {
			return nil, err
		}

		habits = append(habits, habit)
	}

	return habits, nil
}

func (r Repository) MarkHabitChecked(habit models.Habit) error {
	collectionName := "habits"
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(habit.ID)
	if err != nil {
		return fmt.Errorf("could not convert string to ObjectID: %v", err)
	}

	filter := map[string]interface{}{
		"_id": objID,
	}
	exists, err := r.DB.Exists(ctx, collectionName, filter)
	if err != nil {
		return err
	}
	if !exists {
		return ErrHabitNotFound
	}

	err = r.DB.UpdateHabitTracking(ctx, collectionName, habit.ID, time.Now(), true)
	if err != nil {
		return err
	}
	return nil
}
