package main

import (
	"github.com/deenikarim/bookings/pkg/config"
	"github.com/deenikarim/bookings/pkg/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	/**/ //installing middleware
	mux.Use(middleware.Recoverer)

	/**/ //how to use the middleware created
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	/**/ //here is where to set up the routes

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/reservation", handlers.Repo.Reservation)

	/**/
	return mux
}
