package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maxzhirnov/habits/internal/config"
	"github.com/maxzhirnov/habits/internal/database"
	"github.com/maxzhirnov/habits/internal/handlers"
	"github.com/maxzhirnov/habits/internal/repos"
	"github.com/maxzhirnov/habits/internal/service"
	log "github.com/sirupsen/logrus"
	"time"
)

func main() {
	time.Local = time.UTC
	log.Infof("Starting. Local time: %s", time.Now())
	cfg := config.New()
	cfg.Parse()

	mongodb := database.NewMongoConnection(cfg.Server.MongoDBConnection, "habits")
	if err := mongodb.Connect(); err != nil {
		log.Fatal(err)
	}
	defer mongodb.Disconnect(context.Background())

	repo := repos.New(mongodb)
	app := service.NewApp(repo)
	h := handlers.New(app)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/add-new-habit", h.AddNewHabitHandler)
	e.GET("/habits", h.ListHabits)
	e.POST("/mark", h.MarkHabitCheckedForToday)

	e.Logger.Fatal(e.Start(":8080"))
}
