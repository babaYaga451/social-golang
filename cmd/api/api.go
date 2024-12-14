package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/babaYaga451/social/docs"
	"github.com/babaYaga451/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type application struct {
	conf   config
	store  store.Storage
	logger *zap.SugaredLogger
}

type config struct {
	addr   string
	db     dbConfig
	apiURL string
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		docsUrl := fmt.Sprintf("%s/swagger/doc.json", app.conf.addr)
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsUrl), //The url pointing to API definition
		))
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postId}", func(r chi.Router) {
				r.Use(app.postsContextMiddleWare)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.updatePostHandler)
				r.Post("/comments", app.createCommentHandler)
			})
		})

		r.Route("/users", func(r chi.Router) {
			r.Route("/{userId}", func(r chi.Router) {
				r.Use(app.userContextMiddleWare)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	// Docs
	docs.SwaggerInfo.Version = "0.0.1"
	docs.SwaggerInfo.Host = app.conf.apiURL

	srv := &http.Server{
		Addr:         app.conf.addr,
		Handler:      mux,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server started", "addr", app.conf.addr)

	return srv.ListenAndServe()
}
