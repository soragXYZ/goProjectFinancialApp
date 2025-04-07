package config

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v11"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// Global env var spreads across every package
// Not ideal, should be modified later
// https://www.alexedwards.net/blog/organising-database-access
var DB *sql.DB
var Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.DateTime}).With().Timestamp().Logger()
var Conf ConfStruct

type ConfServer struct {
	Port         int           `env:"SERVER_PORT,required" envDefault:"8080"`
	TimeoutRead  time.Duration `env:"SERVER_TIMEOUT_READ,required"`
	TimeoutWrite time.Duration `env:"SERVER_TIMEOUT_WRITE,required"`
	TimeoutIdle  time.Duration `env:"SERVER_TIMEOUT_IDLE,required"`
	LogLevel     string        `env:"SERVER_LOG_LEVEL,required"`
}

type ConfDB struct {
	Host     string `env:"DB_HOST,required"`
	Port     int    `env:"DB_PORT,required"`
	DBName   string `env:"DB_NAME,required"`
	Username string `env:"DB_USER,required"`
	Password string `env:"DB_PASS,required"`
}

type ConfStruct struct {
	Server ConfServer
	DB     ConfDB
}

func Init() {

	// Read env values from .env. Remove this part if your envs are exported from somewhere else
	err := godotenv.Load()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Failed to load .env file")
	}
	// End of the part to remove

	if err := env.Parse(&Conf.DB); err != nil {
		Logger.Fatal().Err(err).Msg("Failed to load env for DB")
	}
	if err := env.Parse(&Conf.Server); err != nil {
		Logger.Fatal().Err(err).Msg("Failed to load env for server")
	}

	// Set log level
	switch Conf.Server.LogLevel {
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		Logger.Fatal().Msg("sdsdsd")
	}

	// Capture connection properties and connect to DB.
	cfg := mysql.Config{
		User:   Conf.DB.Username,
		Passwd: Conf.DB.Password,
		Net:    "tcp",
		Addr:   fmt.Sprintf("%s:%d", Conf.DB.Host, Conf.DB.Port),
		DBName: Conf.DB.DBName,
	}

	DB, err = sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		Logger.Fatal().Err(err).Msg("")
	}

	err = DB.Ping()
	if err != nil {
		Logger.Fatal().Err(err).Msg("")
	}

	Logger.Info().Msg("Successfully ping DB")

}
