package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi"
	"github.com/padurean/purest/internal/database"
	"github.com/padurean/purest/internal/logging"
)

// Start ...
func Start(port string, logger *logging.Logger, db *database.DB) {
	logger.Info().Msg("starting server ...")

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)

	server := newServer(port, logger, db)
	go gracefullShutdown(server, logger, quit, done)

	logger.Info().Msgf("server is ready to handle HTTP requests on port %s", port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatal().Err(err).Msgf("server startup on port %s failed", port)
	}

	<-done
	logger.Info().Msg("server stopped")
}

func gracefullShutdown(server *http.Server, logger *logging.Logger, quit <-chan os.Signal, done chan<- bool) {
	<-quit
	logger.Info().Msg("server is shutting down ...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Msgf("could not gracefully shutdown the server: %v\n", err)
	}
	close(done)
}

func newServer(port string, logger *logging.Logger, db *database.DB) *http.Server {
	router := Router{Router: chi.NewRouter()}
	router.Setup(db, logger)

	return &http.Server{
		Addr:     ":" + port,
		Handler:  router,
		ErrorLog: log.New(logger.Logger, "server: ", 0),
	}
}
