package config

import (
	apiclients "github.com/a-novel/go-api-clients"
	"github.com/rs/zerolog"
	"net/url"
)

func GetAuthClient(logger zerolog.Logger) apiclients.AuthClient {
	authURL, err := new(url.URL).Parse(API.External.AuthAPI)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not parse auth API URL")
	}

	return apiclients.NewAuthClient(authURL)
}

func GetForumClient(logger zerolog.Logger) apiclients.ForumClient {
	forumURL, err := new(url.URL).Parse(API.External.ForumAPI)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not parse forum API URL")
	}

	return apiclients.NewForumClient(forumURL)
}

func GetPermissionsClient(logger zerolog.Logger) apiclients.PermissionsClient {
	permissionsURL, err := new(url.URL).Parse(API.External.PermissionsAPI)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not parse auth API URL")
	}

	return apiclients.NewPermissionsClient(permissionsURL)
}
