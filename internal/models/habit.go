package models

import "time"

type Habit struct {
	ID       string             `json:"id" bson:"_id"`
	UserID   string             `json:"user_id" bson:"user_id"`
	Name     string             `json:"name" bson:"name"`
	Tracking map[time.Time]bool `json:"tracking" bson:"tracking"`
}
