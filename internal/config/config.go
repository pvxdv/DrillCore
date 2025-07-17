package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	env   = "ENV"
	debug = "DEBUG"

	dbHost     = "DB_HOST"
	dbPort     = "DB_PORT"
	dbUser     = "DB_USER"
	dbPassword = "DB_PASS"
	dbName     = "DB_NAME"

	tgToken     = "TG_TOKEN"
	tgBaseURL   = "TG_BASE_URL"
	tgBatchSize = "TG_BATCH_SIZE"
)

var (
	ErrEnvNotExists  = errors.New("environment variable not exists")
	ErrEnvNotCorrect = errors.New("invalid environment variable")
)

type ServiceConfig struct {
	AppEnvs      *AppEnvs
	DbEnvs       *DbEnvs
	TelegramEnvs *TelegramEnvs
}

type AppEnvs struct {
	Env       string
	DebugFlag bool
}

type DbEnvs struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type TelegramEnvs struct {
	Token     string
	BaseUrl   string
	BatchSize int
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
	debugStr, ok := os.LookupEnv(debug)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, debug)
	}

	df, err := strconv.ParseBool(debugStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotCorrect, debugStr)
	}

	e, ok := os.LookupEnv(env)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, env)
	}

	return &AppEnvs{DebugFlag: df, Env: e}, nil
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

	return &DbEnvs{Port: port, Host: host, User: user, Pass: pass, Name: name}, nil
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

	bSizeStr, ok := os.LookupEnv(tgBatchSize)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotExists, tgBatchSize)
	}

	bSize, err := strconv.Atoi(bSizeStr)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrEnvNotCorrect, tgBatchSize)
	}

	return &TelegramEnvs{Token: token, BaseUrl: bUrl, BatchSize: bSize}, nil
}
