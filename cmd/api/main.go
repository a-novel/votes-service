package main

import (
	"fmt"
	"github.com/a-novel/votes-service/config"
	"github.com/a-novel/votes-service/pkg/adapters"
	"github.com/a-novel/votes-service/pkg/dao"
	"github.com/a-novel/votes-service/pkg/handlers"
	"github.com/a-novel/votes-service/pkg/models"
	"github.com/a-novel/votes-service/pkg/services"
)

func main() {
	logger := config.GetLogger()
	authClient := config.GetAuthClient()
	forumClient := config.GetForumClient()

	postgres, closer := config.GetPostgres(logger)
	defer closer()

	votesDAO := dao.NewVotesRepository(postgres)

	votesClients := map[string]models.CheckVoteClient{
		"improveRequest":    adapters.NewImproveRequestVoteClient(forumClient),
		"improveSuggestion": adapters.NewImproveSuggestionVoteClient(forumClient),
	}

	castVoteService := services.NewCastVoteService(votesDAO, authClient, votesClients)
	getUserVoteService := services.NewGetUserVoteService(votesDAO, authClient)
	getVotesSummaryService := services.NewGetVotesSummaryService(votesDAO)
	listUserVotesService := services.NewListUserVotesService(votesDAO, authClient)

	pingHandler := handlers.NewPingHandler()
	healthCheckHandler := handlers.NewHealthCheckHandler(postgres, authClient, forumClient)
	castVoteHandler := handlers.NewCastVoteHandler(castVoteService)
	getUserVoteHandler := handlers.NewGetUserVoteHandler(getUserVoteService)
	getVotesSummaryHandler := handlers.NewGetVotesSummaryHandler(getVotesSummaryService)
	listUserVotesHandler := handlers.NewListUserVotesHandler(listUserVotesService)

	router := config.GetRouter(logger)

	router.GET("/ping", pingHandler.Handle)
	router.GET("/healthcheck", healthCheckHandler.Handle)
	router.POST("/vote", castVoteHandler.Handle)
	router.GET("/vote", getUserVoteHandler.Handle)
	router.GET("/votes/post", getVotesSummaryHandler.Handle)
	router.GET("/votes/user", listUserVotesHandler.Handle)

	if err := router.Run(fmt.Sprintf(":%d", config.API.Port)); err != nil {
		logger.Fatal().Err(err).Msg("a fatal error occurred while running the API, and the server had to shut down")
	}
}
