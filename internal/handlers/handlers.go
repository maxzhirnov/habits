package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	log "github.com/sirupsen/logrus"

	"github.com/maxzhirnov/habits/internal/models"
	"github.com/maxzhirnov/habits/internal/repos"
)

type ApplicationService interface {
	AddNewHabit(habit models.Habit) error
}

type Handlers struct {
	ApplicationService ApplicationService
}

func New(app ApplicationService) *Handlers {
	return &Handlers{
		ApplicationService: app,
	}
}

type HabitDTO struct {
	Name   string `json:"name"`
	UserID string `json:"user_id"`
}

func (h *Handlers) AddNewHabitHandler(c echo.Context) error {
	habit := models.Habit{}
	dto := HabitDTO{}
	if err := c.Bind(&dto); err != nil {
		log.Warn(err)
		return echo.NewHTTPError(http.StatusBadRequest, "error binding json")
	}

	if dto.Name == "" || dto.UserID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name and user_id should be provided")
	}

	habit.Name = dto.Name
	habit.UserID = dto.UserID
	if err := h.ApplicationService.AddNewHabit(habit); err != nil {
		log.Warn(err)
		if errors.Is(err, repos.ErrHabitExists) {
			return echo.NewHTTPError(http.StatusConflict, "habit already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success", "habitName": habit.Name, "userID": dto.UserID})
}
