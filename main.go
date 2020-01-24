package main

import (
	"log"

	"github.com/padurean/purest/auth"
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

	logger.Info().Msg("generating or loading access keys ...")
	if err := auth.GenerateOrLoadKeys(); err != nil {
		logger.Fatal().Err(err).Msgf("error generating or loading access keys")
	}

	//---> TODO OGG: remove these - they are just to test
	var testUserID int64 = 123
	token, err := auth.GenerateToken(testUserID)
	logger.Debug().Msgf("Token: %s, Error: %v", token, err)
	jsonToken, err := auth.VerifyToken(token)
	if err != nil {
		logger.Debug().Msgf("Token verify error: %v", err)
	} else {
		logger.Debug().Msgf("JSONToken: expiration: %s, userID: %d", jsonToken.Expiration, jsonToken.UserID)
	}
	//<---

	logger.Info().Msg("connecting to database ...")
	db := database.MustConnect(env.GetDbDriver(), env.GetDbURL())
	logger.Info().Msg("migrating database ...")
	database.Migrate(db)

	server.Start(env.GetHTTPPort(), logger, db)
}
