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
	DBConfig         DBConfig
	ServerConfig     ServerConfig
	HTTPClientConfig HTTPClientConfig
}

type DBConfig struct {
	URL string
}

type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

type HTTPClientConfig struct {
	AgeBaseURL         string
	GenderBaseURL      string
	NationalityBaseURL string
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
		return nil, fmt.Errorf("invalid read timeout: %w", err)
	}

	writeTimeout, err := strconv.Atoi(os.Getenv("SERVER_WRITE_TIMEOUT"))
	if err != nil {
		return nil, fmt.Errorf("invalid write timeout: %w", err)
	}

	ageURL := os.Getenv("AGE_BASE_URL")
	if ageURL == "" {
		return nil, errors.New("age api URL is empty")
	}

	genderURL := os.Getenv("GENDER_BASE_URL")
	if genderURL == "" {
		return nil, errors.New("gender api URL is empty")
	}

	nationalityURL := os.Getenv("NATIONALITY_BASE_URL")
	if nationalityURL == "" {
		return nil, errors.New("nationality api URL is empty")
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
		HTTPClientConfig: HTTPClientConfig{
			AgeBaseURL:         ageURL,
			GenderBaseURL:      genderURL,
			NationalityBaseURL: nationalityURL,
		},
	}, nil
}
