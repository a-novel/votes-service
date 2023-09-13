package config

import (
	"github.com/rs/zerolog"
	"os"
)

func GetLogger() zerolog.Logger {
	logger := zerolog.New(os.Stdout).
		With().
		Dict(
			"application", zerolog.Dict().
				Str("name", App.Name).
				Str("env", ENV).
				Dict(
					"cors", zerolog.Dict().
						Strs("allowed_origins", App.Frontend.URLs).
						Strs("allowed_methods", Cors.AllowMethods).
						Strs("allowed_headers", Cors.AllowHeaders).
						Strs("exposed_headers", Cors.ExposeHeaders).
						Bool("allow_credentials", Cors.AllowCredentials).
						Dur("max_age", Cors.MaxAge),
				).
				Str("host", API.Host).
				Int("port", API.Port),
		).
		Logger()

	switch ENV {
	case ProdENV:
		logger = logger.With().Timestamp().Logger()
	default:
		logger = logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return logger
}
