package config

import (
	"errors"
	"fmt"
	"os"
)

const (
	debug = "DEBUG"

	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASSWORD"
	dbName     = "DB_NAME"

	tgToken   = "TG_TOKEN"
	tgBaseURL = "TG_BASE_URL"
)

var (
	ErrEnvNotExists = errors.New("environment variable not exists")
)

type ServiceConfig struct {
	AppEnvs      *AppEnvs
	DbEnvs       *DbEnvs
	TelegramEnvs *TelegramEnvs
}

type AppEnvs struct {
	DebugFlag string
}

type DbEnvs struct {
	DbHost string
	DbPort string
	DbUser string
	DbPass string
	DbName string
}

type TelegramEnvs struct {
	Token   string
	BaseUrl string
}

func New() (*ServiceConfig, error) {
	app, err := appEnvs()
	if err != nil {
		return nil, err
	}

	db, err := dbEnvsEnvs()
	if err != nil {
		return nil, err
	}

	tg, err := tgEnvs()
	if err != nil {
		return nil, err
	}

	return &ServiceConfig{
		AppEnvs:      app,
		DbEnvs:       db,
		TelegramEnvs: tg,
	}, nil
}

func appEnvs() (*AppEnvs, error) {
	debugFlag, ok := os.LookupEnv(debug)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, debug)
	}

	return &AppEnvs{DebugFlag: debugFlag}, nil
}

func dbEnvsEnvs() (*DbEnvs, error) {
	host, ok := os.LookupEnv(dbHost)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, dbHost)
	}

	port, ok := os.LookupEnv(dbPort)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, dbPort)
	}

	user, ok := os.LookupEnv(dbUser)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, dbUser)
	}

	pass, ok := os.LookupEnv(dbPassword)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, dbPassword)
	}

	name, ok := os.LookupEnv(dbName)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, dbName)
	}

	return &DbEnvs{DbPort: port, DbHost: host, DbUser: user, DbPass: pass, DbName: name}, nil
}

func tgEnvs() (*TelegramEnvs, error) {
	token, ok := os.LookupEnv(tgToken)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, tgToken)
	}

	bUrl, ok := os.LookupEnv(tgBaseURL)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, tgBaseURL)
	}

	return &TelegramEnvs{Token: token, BaseUrl: bUrl}, nil
}
