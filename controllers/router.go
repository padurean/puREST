package controllers

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/padurean/purest/database"
	"github.com/padurean/purest/env"
	"github.com/padurean/purest/logging"
	"github.com/rs/zerolog/hlog"

	// init Swagger API Docs
	_ "github.com/padurean/purest/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router ...
type Router struct {
	chi.Router
}

func (router Router) setupMiddlewares(db *database.DB, logger *logging.Logger) {
	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Recoverer,
		middleware.Heartbeat("/ping"),
		// Set a timeout value on the request context (ctx), that will signal
		// through ctx.Done() that the request has timed out and further
		// processing should be stopped.
		middleware.Timeout(60*time.Second),
		// set DB connection on request context
		middleware.WithValue("db", db),

		//--> logging middleware
		hlog.NewHandler(*logger.Logger),
		hlog.RemoteAddrHandler("ip"),
		hlog.UserAgentHandler("user_agent"),
		hlog.RefererHandler("referer"),
		hlog.RequestIDHandler("req_id", "X-Request-Id"),
		//<--

		// set the logger on the request context
		middleware.WithValue("logger", logger),
	)

	if env.GetLogRequests() {
		router.Use(
			hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
				hlog.FromRequest(r).Info().
					Str("method", r.Method).
					Str("url", r.URL.String()).
					Int("status", status).
					Int("size", size).
					Str("duration", duration.String()).
					Msg("")
			}))
	}
}

func (router Router) setupRoutes() {
	// setup API Docs (Swagger) routes
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		logging.Simple(r).Debug().Msg("Simple Logger: Buona Sera, Siniora!!! >:D")
		logging.Detailed(r).Debug().Msg("Detailed Logger: Buona Sera, Siniora!!! >:D")
		w.Write([]byte("Hello, my name is puREST. Pleased to meet you! :)"))
	})
}

// Setup ...
func (router Router) Setup(db *database.DB, logger *logging.Logger) {
	router.setupMiddlewares(db, logger)
	router.setupRoutes()
}
