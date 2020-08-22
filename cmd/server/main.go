package main

import (
	"log"

	"github.com/padurean/purest/internal/auth"
	"github.com/padurean/purest/internal/database"
	"github.com/padurean/purest/internal/env"
	"github.com/padurean/purest/internal/logging"
	"github.com/padurean/purest/internal/server"
)

// @title puREST API
// @version 1.0
// @description Golang REST API boilerplate with authentication using PASETO tokens, RBAC authorization, PostgreSQL and Swagger for API docs.
// @license.name MIT
// @license.url https://tldrlegal.com/license/mit-license
// @basePath /api/v1
func main() {
	logger := logging.FromEnv()
	// also set the logger as output for any standard log usage
	log.SetFlags(0)
	log.SetOutput(logger.Logger)

	logger.Info().Msg("generating or loading access keys ...")
	if err := auth.GenerateOrLoadKeys(); err != nil {
		logger.Fatal().Err(err).Msgf("error generating or loading access keys")
	}

	logger.Info().Msg("connecting to database ...")
	db := database.MustConnect(env.GetDbDriver(), env.GetDbURL())
	logger.Info().Msg("migrating database ...")
	database.Migrate(db)

	server.Start(env.GetHTTPPort(), logger, db)
}
