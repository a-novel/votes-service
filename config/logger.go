package config

import (
	"github.com/rs/zerolog"
	"os"
)

func GetLogger() zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		With().
		Dict("application", zerolog.Dict().Str("name", App.Name).Str("env", ENV)).
		Logger()

	switch ENV {
	case ProdENV:
		logger = logger.With().Timestamp().Logger()
	default:
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return logger
}
