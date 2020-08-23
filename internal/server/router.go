package server

import (
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/docgen"
	icontext "github.com/padurean/purest/internal/context"
	"github.com/padurean/purest/internal/controller"
	"github.com/padurean/purest/internal/database"
	"github.com/padurean/purest/internal/env"
	"github.com/padurean/purest/internal/logging"
	"github.com/rs/zerolog/hlog"

	// init Swagger API Docs
	_ "github.com/padurean/purest/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Router ...
type Router struct {
	chi.Router
}

func (router Router) setupCommonMiddlewares(db *database.DB, logger *logging.Logger) {
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
		middleware.WithValue(icontext.KeyDB, db),

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
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		msg := "Hello, my name is puREST. Pleased to meet you! :)"
		logging.Simple(r).Debug().Msgf("Simple Logger: %s", msg)
		logging.Detailed(r).Debug().Msgf("Detailed Logger: %s", msg)
		w.Write([]byte(msg))
	})

	// setup API Docs (Swagger) routes
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8000/swagger/doc.json"),
	))

	router.Route("/api", func(router chi.Router) {
		router.Route("/v1", func(router chi.Router) {

			router.Route("/users", func(router chi.Router) {
				router.With(controller.UserCtx).Post("/sign-in/{usernameOrEmail}", controller.UserSignIn)
				routerWithAuth := router.With(authenticate)
				routerWithAuth.Post("/", controller.UserCreate)
				routerWithAuth.With(paginate).Get("/", controller.UserList)
				routerWithAuth.Route("/{id}", func(router chi.Router) {
					router.Use(controller.UserCtx)
					router.Get("/", controller.UserGet)
					router.Put("/", controller.UserUpdate)
					router.Delete("/", controller.UserDelete)
				})
			})

		})
	})
}

func (router Router) generateAPIDocs(logger *logging.Logger) {
	// jsonDocs := docgen.JSONRoutesDoc(router)
	mdDocs := docgen.MarkdownRoutesDoc(router, docgen.MarkdownOpts{
		ProjectPath: "github.com/padurean/purest",
		Intro:       "Welcome to the puREST generated docs!",
	})
	mdDocsFilename := "API-docs.md"
	mdDocsFile, err := os.Create(mdDocsFilename)
	if err != nil {
		logger.Err(err).Msgf("error creating API docs markdown file %s", mdDocsFilename)
	}
	defer mdDocsFile.Close()
	_, err = mdDocsFile.WriteString(mdDocs)
	if err != nil {
		logger.Err(err).Msgf("error writing markdown content to API docs file %s", mdDocsFilename)
	}
	logger.Info().Msgf("%s API docs file written successfully", mdDocsFilename)
}

// Setup ...
func (router Router) Setup(db *database.DB, logger *logging.Logger) {
	router.setupCommonMiddlewares(db, logger)
	router.setupRoutes()
	if env.GetAppEnv() == env.Development {
		router.generateAPIDocs(logger)
	}
}
