package config

import (
	"errors"
	"flag"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

var Config struct {
	TelegramToken string `env:"TELEGRAM_TOKEN"`
	DatabaseURI   string `env:"DATABASE_URI"`

	MigrationsFlag bool `env:"MIGRATIONS_FLAG"`
	ProductionMode bool `env:"PRODUCTION_MODE"`
}

func ParseConfig() error {
	flag.StringVar(&Config.TelegramToken, "token", "", "telegram bot token")
	flag.StringVar(&Config.DatabaseURI, "database", "", "uri of database")
	flag.BoolVar(&Config.MigrationsFlag, "migrations", true, "need migrations?")
	flag.BoolVar(&Config.ProductionMode, "production", false, "production mode?")
	flag.Parse()

	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}

	return env.Parse(&Config)
}
