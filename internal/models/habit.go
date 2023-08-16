package models

import (
	"github.com/maxzhirnov/habits/internal/utils"
	"time"
)

type Habit struct {
	ID       string             `json:"id" bson:"_id,omitempty"`
	UserID   string             `json:"user_id" bson:"user_id"`
	Name     string             `json:"name" bson:"name"`
	CratedAt time.Time          `json:"created" bson:"created"`
	Streak   int                `json:"streak" bson:"streak"`
	Tracking map[time.Time]bool `json:"tracking" bson:"tracking"`
}

func (h *Habit) IsDoneToday() bool {
	today := utils.DateOnly(time.Now())
	done, exists := h.Tracking[today]
	return exists && done
}
