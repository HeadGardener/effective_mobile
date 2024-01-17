package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

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
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
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

	readTimeout, err := strconv.Atoi(os.Getenv("SERVER_READ_TIMEOUT"))
	if err != nil {
		return nil, errors.Join(errors.New("invalid read timeout"), err)
	}

	writeTimeout, err := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		return nil, errors.Join(errors.New("invalid write timeout"), err)
	}

	return &Config{
		DBConfig: DBConfig{
			URL: dburl,
		},
		ServerConfig: ServerConfig{
			Port:         srvport,
			ReadTimeout:  time.Duration(readTimeout),
			WriteTimeout: time.Duration(writeTimeout),
		},
	}, nil
}
