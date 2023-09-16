package main

import (
	"context"
	"fmt"
	"github.com/a-novel/bunovel"
	"github.com/a-novel/go-apis"
	"github.com/a-novel/votes-service/config"
	"github.com/a-novel/votes-service/migrations"
	"github.com/a-novel/votes-service/pkg/adapters"
	"github.com/a-novel/votes-service/pkg/dao"
	"github.com/a-novel/votes-service/pkg/handlers"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/services"
	"io/fs"
)

func main() {
	ctx := context.Background()
	logger := config.GetLogger()
	authClient := config.GetAuthClient(logger)
	forumClient := config.GetForumClient(logger)
	permissionsClient := config.GetPermissionsClient(logger)

	postgres, sql, err := bunovel.NewClient(ctx, bunovel.Config{
		Driver:                &bunovel.PGDriver{DSN: config.Postgres.DSN, AppName: config.App.Name},
		Migrations:            &bunovel.MigrateConfig{Files: []fs.FS{migrations.Migrations}},
		DiscardUnknownColumns: true,
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("error connecting to postgres")
	}
	defer func() {
		_ = postgres.Close()
		_ = sql.Close()
	}()

	votesDAO := dao.NewVotesRepository(postgres)

	votesClients := map[string]models.CheckVoteClient{
		"improveRequest":    adapters.NewImproveRequestVoteClient(forumClient, permissionsClient),
		"improveSuggestion": adapters.NewImproveSuggestionVoteClient(forumClient, permissionsClient),
	}

	castVoteService := services.NewCastVoteService(votesDAO, authClient, votesClients)
	getUserVoteService := services.NewGetUserVoteService(votesDAO, authClient)
	getVotesSummaryService := services.NewGetVotesSummaryService(votesDAO)
	listUserVotesService := services.NewListUserVotesService(votesDAO, authClient)

	castVoteHandler := handlers.NewCastVoteHandler(castVoteService)
	getUserVoteHandler := handlers.NewGetUserVoteHandler(getUserVoteService)
	getVotesSummaryHandler := handlers.NewGetVotesSummaryHandler(getVotesSummaryService)
	listUserVotesHandler := handlers.NewListUserVotesHandler(listUserVotesService)

	router := apis.GetRouter(apis.RouterConfig{
		Logger:    logger,
		ProjectID: config.Deploy.ProjectID,
		CORS:      apis.GetCORS(config.App.Frontend.URLs),
		Prod:      config.ENV == config.ProdENV,
		Health: map[string]apis.HealthChecker{
			"postgres": func() error {
				return postgres.PingContext(ctx)
			},
			"auth-client": func() error {
				return authClient.Ping(ctx)
			},
			"forum-client": func() error {
				return forumClient.Ping(ctx)
			},
			"permissions-client": func() error {
				return permissionsClient.Ping(ctx)
			},
		},
	})

	router.POST("/vote", castVoteHandler.Handle)
	router.GET("/vote", getUserVoteHandler.Handle)
	router.GET("/votes/post", getVotesSummaryHandler.Handle)
	router.GET("/votes/user", listUserVotesHandler.Handle)

	if err := router.Run(fmt.Sprintf(":%d", config.API.Port)); err != nil {
		logger.Fatal().Err(err).Msg("a fatal error occurred while running the API, and the server had to shut down")
	}
}
