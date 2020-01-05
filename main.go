package main

import (
	"log"

	"github.com/padurean/purest/database"
	"github.com/padurean/purest/env"
	"github.com/padurean/purest/logging"
	"github.com/padurean/purest/server"
)

// @title puREST API
// @version 1.0
// @description Golang REST API boilerplate with JWT-based authentication, RBAC authorization and PostgreSQL.
// @license.name MIT
// @license.url https://tldrlegal.com/license/mit-license
// @host purecore.ro/purest
// @BasePath /api
func main() {
	logger := logging.FromEnv()
	// also set the logger as output for any standard log usage
	log.SetFlags(0)
	log.SetOutput(logger.Logger)

	logger.Info().Msg("connecting to database ...")
	db := database.MustConnect(env.GetDbDriver(), env.GetDbURL())
	logger.Info().Msg("migrating database ...")
	database.Migrate(db)

	server.Start(env.GetHTTPPort(), logger, db)
}
