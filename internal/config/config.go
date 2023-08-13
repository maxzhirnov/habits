package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Warning(err)
	}
}

type Config struct {
	Environment       string
	MongoDBConnection string
}

func New() *Config {
	return &Config{
		Environment: os.Getenv("ENV"),
	}
}

func (c *Config) Parse() {
	var ok bool
	c.MongoDBConnection, ok = os.LookupEnv("MONGO_CONN")
	if !ok {
		log.Warn("Couldn't parse mongodb connection string from environment variables")
	}
	log.Info("Configs parsed")
}
