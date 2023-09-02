package config

import (
	"fmt"
	"os"

	"github/wry-0313/exchange/validator.go"

	"github.com/joho/godotenv"
)

const (
	keyDBHost     = "DB_HOST"
	keyDBPort     = "DB_PORT"
	keyDBName     = "DB_NAME"
	keyDBUser     = "DB_USER"
	keyDBPassword = "DB_PASSWORD"

	keyInternalNetwork = "INTERNAL_NETWORK"

	keyEnv        = "ENV"
	keyServerPort = "SERVER_PORT"

	valEnvDev = "DEVELOPMENT"
)

type Config struct {
	DB         DatabaseConfig
	ServerPort string
}

func Load(file string) (*Config, error) {
	env := os.Getenv(keyEnv)
	if env == valEnvDev {
		// Load .env file if in development
		err := godotenv.Load(file)
		if err != nil {
			return nil, fmt.Errorf("error loading .env file: %w", err)
		}
	}

	databaseConfig, err := getDatabaseConfig()
	if err != nil {
		return nil, err
	}

	serverPort := os.Getenv(keyServerPort)

	return &Config{
		DB:         databaseConfig,
		ServerPort: serverPort,
	}, nil
}

// DatabaseConfig encapsulates all the config values for the database.
type DatabaseConfig struct {
	Host     string `validate:"required"`
	Port     string `validate:"required"`
	Name     string `validate:"required"`
	User     string `validate:"required"`
	Password string `validate:"required"`
}

// Validate checks that all values are properly loaded into the database config.
func (dbConfig *DatabaseConfig) Validate() error {
	validate := validator.New()
	if err := validate.Struct(dbConfig); err != nil {
		return fmt.Errorf("missing database env var: %v", err)
	}
	return nil
}

func getDatabaseConfig() (DatabaseConfig, error) {
	databaseConfig := DatabaseConfig{
		Host:     os.Getenv(keyDBHost),
		Port:     os.Getenv(keyDBPort),
		Name:     os.Getenv(keyDBName),
		User:     os.Getenv(keyDBUser),
		Password: os.Getenv(keyDBPassword),
	}

	// This allows running tests from outside the docker network assuming your local
	// development environment has ports exposed
	if os.Getenv(keyInternalNetwork) == "false" {
		databaseConfig.Host = "localhost"
	}

	// validate all db params are available
	if err := databaseConfig.Validate(); err != nil {
		return DatabaseConfig{}, err
	}

	return databaseConfig, nil
}
