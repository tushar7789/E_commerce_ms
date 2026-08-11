package main

import (
	"log"
	"net/http"
	"time"

	domain "github.com/tushar7789/E_commerce_ms/internal/domain"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *domain.application) mount() http.Handler {
	jwt_cfg := Load()
	r := chi.NewRouter()

	// Midlewares
	r.Use(middleware.RequestID) // important for rate limiting (ddos)
  	r.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
  	r.Use(middleware.Logger)
  	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("all good!"))
	})

	authService := auth.NewService(db.New(app.db), jwt_cfg.AccessSecret, jwt_cfg.RefreshSecret, jwt_cfg.AccessTTL, jwt_cfg.RefreshTTL)
	authHandler := auth.NewHandler(authService)
	r.Post("/users", authHandler.CreateUser)
	r.Post("/login", authHandler.Login)
	r.Post("/refresh", authHandler.Refresh)

	return r
}

func (app *domain.application) run(h http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server has started at addr %s", app.config.addr)

	return srv.ListenAndServe()
}
