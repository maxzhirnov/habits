package models

type User struct {
	ID     string  `json:"id" bson:"id"`
	Name   string  `json:"name" bson:"name"`
	Habits []Habit `json:"habits" bson:"habits"`
}
