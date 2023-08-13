package main

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/maxzhirnov/habits/internal/config"
	"github.com/maxzhirnov/habits/internal/database"
	"github.com/maxzhirnov/habits/internal/handlers"
	"github.com/maxzhirnov/habits/internal/repos"
	"github.com/maxzhirnov/habits/internal/services"
	log "github.com/sirupsen/logrus"
)

func main() {
	cfg := config.New()
	cfg.Parse()

	mongodb := database.NewMongoConnection(cfg.MongoDBConnection, "habits")
	if err := mongodb.Connect(); err != nil {
		log.Fatal(err)
	}
	defer mongodb.Disconnect(context.Background())

	repo := repos.New(mongodb)
	app := services.NewAppService(repo)
	h := handlers.New(app)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/add-new-habit", h.AddNewHabitHandler)
	e.Logger.Fatal(e.Start(":8080"))

}
