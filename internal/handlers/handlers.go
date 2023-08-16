package handlers

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/maxzhirnov/habits/internal/models"
	"github.com/maxzhirnov/habits/internal/repos"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

type ApplicationService interface {
	AddNewHabit(habit models.Habit) error
	GetAllUserHabits(userID string) ([]models.Habit, error)
	MarkHabitChecked(h models.Habit) error
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
	ID           string `json:"id"`
	Name         string `json:"name"`
	UserID       string `json:"user_id"`
	CreationDate string `json:"created"`
}

func (dto HabitDTO) convertToHabit() (models.Habit, error) {
	const layout = "2006-01-02 15:04:05 -0700 MST"
	parsedTime, err := time.Parse(layout, dto.CreationDate)
	if err != nil {
		log.Warn(err)
	}

	return models.Habit{
		ID:       dto.ID,
		UserID:   dto.UserID,
		Name:     dto.Name,
		CratedAt: parsedTime,
	}, nil
}

func (h *Handlers) AddNewHabitHandler(c echo.Context) error {
	// Читаем только чтобы залогировать боди
	b, err := io.ReadAll(c.Request().Body)
	if err != nil {
		log.Warn(err)
	}
	log.Infof("Got request with body: %s", b)
	// Возвращаем данные обратно в тело запроса
	c.Request().Body = io.NopCloser(bytes.NewBuffer(b))

	dto := HabitDTO{}
	if err := c.Bind(&dto); err != nil {
		log.Warnf("error binding json: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "error binding json")
	}

	if dto.Name == "" || dto.UserID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name and user_id should be provided")
	}

	habit, err := dto.convertToHabit()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if err := h.ApplicationService.AddNewHabit(habit); err != nil {
		log.Warn(err)
		if errors.Is(err, repos.ErrHabitExists) {
			return echo.NewHTTPError(http.StatusConflict, "habit already exists")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "success", "habitName": habit.Name, "userID": dto.UserID})
}

func (h *Handlers) ListHabits(e echo.Context) error {
	params := e.QueryParams()
	userID := params.Get("userid")
	if userID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Errorf("you need to provide userid"))
	}
	habits, err := h.ApplicationService.GetAllUserHabits(userID)
	if err != nil {
		return err
	}

	return e.JSON(http.StatusOK, habits)
}

func (h *Handlers) MarkHabitCheckedForToday(c echo.Context) error {
	dto := HabitDTO{}
	if err := c.Bind(&dto); err != nil {
		log.Warnf("error binding json: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "error binding json")
	}

	if dto.ID == "" || dto.UserID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "name and user_id should be provided")
	}

	habit, err := dto.convertToHabit()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err = h.ApplicationService.MarkHabitChecked(habit)
	if err != nil {
		return err
	}

	return nil
}
