package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig     DBConfig
	ServerConfig ServerConfig
}

type DBConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

func Init(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		return nil, errors.New("db url is empty")
	}

	srvport := os.Getenv("SERVER_PORT")
	if srvport == "" {
		return nil, errors.New("server port is empty")
	}

	return &Config{
		DBConfig: DBConfig{
			URL: dburl,
		},
		ServerConfig: ServerConfig{
			Port: srvport,
		},
	}, nil
}
