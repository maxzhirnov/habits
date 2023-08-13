package models

import "time"

type Habit struct {
	UserID   string             `json:"user_id" bson:"user_id"`
	Name     string             `json:"name" bson:"name"`
	Tracking map[time.Time]bool `json:"tracking" bson:"tracking"`
}
