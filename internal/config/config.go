package config

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Warning(err)
	}
}

type Config struct {
	Server      ServerConfig `yaml:"server"`
	Client      ClientConfig `yaml:"client"`
	Environment string
}

type ClientConfig struct {
	Host          string `yaml:"host"`
	Port          int    `yaml:"port"`
	TelegramToken string
}

type ServerConfig struct {
	MongoDBConnection string
}

func New() *Config {
	return &Config{
		Environment: os.Getenv("ENV"),
	}
}

func (c *Config) Parse() {
	//config yaml
	data, err := os.ReadFile("config/config.yml")
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	if err := yaml.Unmarshal(data, &c); err != nil {
		log.Fatalf("Error parsing config file: %v", err)
	}

	//Env
	var ok bool
	c.Server.MongoDBConnection, ok = os.LookupEnv("MONGO_CONN")
	if !ok {
		log.Warn("Couldn't parse mongodb connection string from environment variables")
	}

	c.Client.TelegramToken, ok = os.LookupEnv("CLIENT_BOT_TOKEN")
	if !ok {
		log.Warn("Couldn't parse telegram token string from environment variables")
	}
	log.Info("Configs parsed")
}
